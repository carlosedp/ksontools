package yaml2jsonnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/bryanl/woowoo/jsonnetutil"
	"github.com/bryanl/woowoo/node"
	"github.com/bryanl/woowoo/params"
	"github.com/google/go-jsonnet/ast"
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
	root, err := jsonnetutil.Import(k8sLib)
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
		componentName := generateComponentName(c.sourceFile)
		doc, err := NewDocument(componentName, r, c.RootNode)
		if err != nil {
			return errors.Wrap(err, "parse document")
		}

		s, err := doc.GenerateComponent()
		if err != nil {
			return errors.Wrap(err, "generate jsonnet")
		}

		if err := doc.UpdateParams(c.updateParams); err != nil {
			return errors.Wrap(err, "update params")
		}

		fmt.Println(s)
	}

	return nil
}

func (c *Conversion) updateParams(componentName string, values map[string]interface{}) error {
	update, err := params.Update([]string{"components", componentName}, paramsSource, values)
	if err != nil {
		return errors.Wrap(err, "update params")
	}

	fmt.Println(update)

	return nil
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

var (
	paramsSource = `{
  global: {
  },
  // Component-level parameters, defined initially from 'ks prototype use ...'
  // Each object below should correspond to a component in the components/ directory
  components: {
    "guestbook-ui": {
      containerPort: 80,
      image: "gcr.io/heptio-images/ks-guestbook-demo:0.1",
      name: "guestbook-ui",
      replicas: 1,
      servicePort: 80,
      type: "ClusterIP",
    },
  },
}
	  `
)
