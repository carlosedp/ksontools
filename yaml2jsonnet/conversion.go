package yaml2jsonnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/bryanl/woowoo/node"
	"github.com/bryanl/woowoo/pkg/docparser"
	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

const (
	docSeparator = "---"
)

// Conversion converts YAML to ksonnet.
type Conversion struct {
	RootNode   ast.Node
	Sources    []io.Reader
	sourceFile string
}

// NewConversion creates a Conversion.
func NewConversion(source, k8sLib string) (*Conversion, error) {
	root, err := ImportJsonnet(k8sLib)
	if err != nil {
		return nil, errors.Wrap(err, "read ksonnet lib")
	}

	node.FindMembers(root)

	readers, err := importSource(source)
	if err != nil {
		return nil, errors.Wrap(err, "import source")
	}

	c := &Conversion{
		RootNode:   root,
		Sources:    readers,
		sourceFile: source,
	}

	return c, nil
}

// Process processes the documents supplied.
func (c *Conversion) Process() error {
	for _, r := range c.Sources {
		doc, err := NewDocument(r, c.RootNode)
		if err != nil {
			return errors.Wrap(err, "parse document")
		}

		componentName := generateComponentName(c.sourceFile)
		s, err := doc.GenerateComponent2(componentName)
		if err != nil {
			return errors.Wrap(err, "generate libsonnet")
		}

		fmt.Println(s)
	}

	return nil
}

// ImportJsonnet imports jsonnet from a path.
func ImportJsonnet(fileName string) (*astext.Object, error) {
	if fileName == "" {
		return nil, errors.New("filename was blank")
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "read lib")
	}

	tokens, err := docparser.Lex(fileName, string(b))
	if err != nil {
		return nil, errors.Wrap(err, "lex lib")
	}

	node, err := docparser.Parse(tokens)
	if err != nil {
		return nil, errors.Wrap(err, "parse lib")
	}

	root, ok := node.(*astext.Object)
	if !ok {
		return nil, errors.New("root was not an object")
	}

	return root, nil
}

func importSource(source string) ([]io.Reader, error) {
	if err := checkSource(source); err != nil {
		return nil, errors.Wrap(err, "check source")
	}

	f, err := os.Open(source)
	if err != nil {
		return nil, errors.Wrap(err, "open source")
	}
	defer f.Close()

	bufs := make([]bytes.Buffer, 1)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if t == docSeparator {
			bufs = append(bufs, bytes.Buffer{})
			continue
		}

		bufs[len(bufs)-1].WriteString(t)
		bufs[len(bufs)-1].WriteByte('\n')
	}

	var readers []io.Reader
	for i := range bufs {
		readers = append(readers, &bufs[i])
	}

	return readers, nil
}

func checkSource(source string) error {
	if source == "" {
		return errors.New("source is empty")
	}

	if _, err := os.Stat(source); err != nil {
		if os.IsNotExist(err) {
			return errors.New("source does not exist")
		}

		return errors.Wrap(err, "could not stat source")
	}

	return nil
}

func generateComponentName(inputFileName string) string {
	inputFileName = filepath.Base(inputFileName)
	componentFile := strings.TrimSuffix(inputFileName, filepath.Ext(inputFileName))
	return strcase.ToLowerCamel(componentFile)
}
