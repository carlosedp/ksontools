package yaml2jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"

	"github.com/go-yaml/yaml"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

type Document struct {
	Properties Properties
	GVK        GVK
	root       *astext.Object
}

func NewDocument(r io.Reader, root ast.Node) (*Document, error) {
	obj, ok := root.(*astext.Object)
	if !ok {
		return nil, errors.New("root is not an *ast.Object")
	}

	doc := &Document{
		Properties: Properties{},
		root:       obj,
	}

	ts, err := importYaml(r, doc.Properties)
	if err != nil {
		return nil, err
	}

	gvk, err := ts.GVK()
	if err != nil {
		return nil, errors.Wrap(err, "type spec is invalid")
	}

	doc.GVK = gvk

	return doc, nil
}

func importYaml(r io.Reader, props Properties) (TypeSpec, error) {
	var m map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, errors.Wrap(err, "decode yaml")
	}

	ts := TypeSpec{}

	for k, v := range m {
		switch k {
		case "kind", "apiVersion":
			ts[k] = v.(string)
		default:
			props[k] = v
		}
	}

	return ts, nil
}

func (d *Document) Selector() string {
	return fmt.Sprintf("k.%s.%s.%s", d.GVK.Group, d.GVK.Version, d.GVK.Kind)
}

func (d *Document) Generate() (string, error) {
	selector := d.Selector()
	comp := NewComponent()

	comp.AddDeclaration(
		Declaration{
			Name:  d.GVK.Kind,
			Value: NewDeclarationApply(selector),
		})

	nn := NewNode("root", d.root)

	locals := NewLocals(d.GVK.Kind)

	spew.Dump(d.GVK)
	paths := d.Properties.Paths(d.GVK)
	for _, path := range paths {
		sr, realPath, err := nn.Search(path.Path...)
		if err != nil {
			return "", errors.Wrapf(err, "search path %s", strings.Join(path.Path, "."))
		}

		manifestPath := realPath[4:]
		var paramName bytes.Buffer
		for i := range manifestPath {
			part := manifestPath[i]
			if i > 0 {
				part = strings.Title(part)
			}

			paramName.WriteString(part)
		}

		v, err := d.Properties.Value(manifestPath)
		if err != nil {
			return "", errors.Wrapf(err, "retrieve manifest values for %s", strings.Join(path.Path, "."))
		}

		if err := comp.AddParam(paramName.String(), v); err != nil {
			return "", errors.Wrapf(err, "add param %s to component", paramName.String())
		}

		k := strings.Join(realPath[:len(realPath)-1], ".")
		entry := LocalEntry{
			Path:      k,
			Setter:    sr.Setter,
			ParamName: paramName.String(),
		}

		locals.Add(entry)
	}

	var kindParts []string

	decls, err := locals.Generate()
	if err != nil {
		return "", errors.Wrap(err, "generate locals")
	}

	for _, decl := range decls {
		kindParts = append(kindParts, decl.Name)
		comp.AddDeclaration(decl)
	}

	nodeInit := fmt.Sprintf("init%s", strings.Title(d.GVK.Kind))

	comp.AddDeclaration(Declaration{
		Name:  nodeInit,
		Value: NewDeclarationNoder(nodemaker.ApplyCall(fmt.Sprintf("%s.new", d.GVK.Kind))),
	})

	n := nodemaker.NewVar(nodeInit)

	var left nodemaker.Noder = n
	for _, name := range kindParts {
		left = nodemaker.NewBinary(left, nodemaker.NewVar(name), nodemaker.BopPlus)
	}

	s, err := comp.Generate(left.Node())
	if err != nil {
		return "", errors.Wrap(err, "generate component")
	}

	return s, nil
}
