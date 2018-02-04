package yaml2jsonnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bryanl/woowoo/yaml2jsonnet/pkg/parser"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/pkg/errors"
)

const (
	docSeparator = "---"
)

type Conversion struct {
	RootNode ast.Node
	Sources  []io.Reader
}

func NewConversion(source, k8sLib string) (*Conversion, error) {
	node, err := ImportJsonnet(k8sLib)
	if err != nil {
		return nil, errors.Wrap(err, "read ksonnet lib")
	}

	readers, err := importSource(source)
	if err != nil {
		return nil, errors.Wrap(err, "import source")
	}

	c := &Conversion{
		RootNode: node,
		Sources:  readers,
	}

	return c, nil
}

func (c *Conversion) Process() error {
	for _, r := range c.Sources {
		doc, err := NewDocument(r, c.RootNode)
		if err != nil {
			return errors.Wrap(err, "parse document")
		}

		s, err := doc.Generate()
		if err != nil {
			return errors.Wrap(err, "generate libsonnet")
		}

		fmt.Println(s)
	}

	return nil
}

func ImportJsonnet(fileName string) (ast.Node, error) {
	if fileName == "" {
		return nil, errors.New("filename was blank")
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "read lib")
	}

	tokens, err := parser.Lex(fileName, string(b))
	if err != nil {
		return nil, errors.Wrap(err, "lex lib")
	}

	node, err := parser.Parse(tokens)
	if err != nil {
		return nil, errors.Wrap(err, "parse lib")
	}

	return node, nil
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
