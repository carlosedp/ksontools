package yaml2jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"

	"github.com/go-yaml/yaml"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

// Document creates a ksonnet document for describing a resource.
type Document struct {
	Properties Properties
	GVK        GVK
	root       *astext.Object
}

// NewDocument creates an instance of Document.
func NewDocument(r io.Reader, root ast.Node) (*Document, error) {
	obj, ok := root.(*astext.Object)
	if !ok {
		return nil, errors.New("root is not an *ast.Object")
	}

	doc := &Document{
		root: obj,
	}

	ts, props, err := importYaml(r)
	if err != nil {
		return nil, err
	}

	doc.Properties = props

	gvk, err := ts.GVK()
	if err != nil {
		return nil, errors.Wrap(err, "type spec is invalid")
	}

	doc.GVK = gvk

	return doc, nil
}

func importYaml(r io.Reader) (TypeSpec, Properties, error) {
	var m map[interface{}]interface{}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, nil, errors.Wrap(err, "decode yaml")
	}

	ts := TypeSpec{}
	props := Properties{}

	for k, v := range m {
		switch k {
		case "kind", "apiVersion":
			ts[k.(string)] = v.(string)
		default:
			props[k] = v
		}
	}

	return ts, props, nil
}

// Selector is the selector for the resource this document represents.
func (d *Document) Selector() string {
	g := d.GVK.Group()

	path := append(g, d.GVK.Version, d.GVK.Kind)
	return fmt.Sprintf("k.%s", strings.Join(path, "."))
}

// Generate generates the document.
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

	paths := d.Properties.Paths(d.GVK)
	for _, path := range paths {
		sr, err := nn.Search(path.Path...)
		if err != nil {
			return "", errors.Wrapf(err, "search path %s", strings.Join(path.Path, "."))
		}

		manifestPath := sr.MatchedPath[4:]
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

		k := strings.Join(sr.MatchedPath[:len(sr.MatchedPath)-1], ".")
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

// GenerateComponent generates a ksonnet component for the document.
func (d *Document) GenerateComponent() (string, error) {

	ctor := nodemaker.NewVar("abcCustomResourceDefinition")
	bodyArgs := nodemaker.NewArray([]nodemaker.Noder{ctor})
	call := nodemaker.NewCall("k.core.v1.list.new")
	object := nodemaker.NewApply(call, []nodemaker.Noder{bodyArgs}, nil)

	resource, err := d.addResource()
	if err != nil {
		return "", err
	}
	resource.Body = object.Node()

	ctor2, err := d.createConstructor()
	if err != nil {
		return "", err
	}

	ctor2.Body = resource

	header := d.declarations(ctor2)

	return d.render(header)

}

type componentImport struct {
	name     string
	location string
}

func (d *Document) declarations(next ast.Node) *ast.Local {
	imports := []componentImport{
		{name: "k", location: "k.libsonnet"},
		{name: "stdlib", location: "stdlib.libsonnet"},
	}

	var declRoot *ast.Local
	var declLast *ast.Local
	for _, imp := range imports {
		decl := Declaration{Name: imp.name, Value: NewDeclarationImport(imp.location)}
		obj := genDeclaration(decl)

		if declRoot == nil {
			declRoot = obj
			declLast = obj
		} else {
			declLast.Body = obj
			declLast = obj
		}
	}

	declLast.Body = next

	return declRoot
}

func (d *Document) createConstructor() (*ast.Local, error) {
	// local createCustomResourceDefinition(params) = {}

	paths, err := d.resolvedPaths()
	if err != nil {
		return nil, err
	}

	spew.Fdump(ioutil.Discard, paths)

	fn := nodemaker.NewFunction([]string{"params"}, nodemaker.NewObject())
	decl := Declaration{
		Name:  "createCustomResourceDefinition",
		Value: NewDeclarationNoder(fn),
	}

	return genDeclaration(decl), nil
}

func (d *Document) addResource() (*ast.Local, error) {
	kind := strings.Title(d.GVK.Kind)

	var resourceName string
	for _, char := range kind {
		if unicode.IsUpper(char) {
			resourceName += string(char)
		}
	}

	resourceName = strings.ToLower(resourceName)

	call := nodemaker.NewCall(fmt.Sprintf("create%s", kind))
	args := nodemaker.NewCall("params")
	apply := nodemaker.NewApply(call, []nodemaker.Noder{args}, nil)

	decl := Declaration{
		Name:  resourceName,
		Value: NewDeclarationNoder(apply),
	}

	return genDeclaration(decl), nil
}

func (d *Document) render(root ast.Node) (string, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, root); err != nil {
		return "", errors.Wrap(err, "create jsonnet")
	}

	return buf.String(), nil
}

type documentValues struct {
	setter string
	value  interface{}
}

func (d *Document) resolvedPaths2() (map[string]documentValues, error) {
	nn := NewNode("root", d.root)

	m := make(map[string]documentValues)
	cache := make(map[string]bool)

	paths := d.Properties.Paths(d.GVK)
	for _, path := range paths {
		item, err := nn.Search2(path.Path...)
		if err != nil {
			continue
		}

		var manifestPath []string
		var found bool
		for _, p := range item.Path {
			if p == d.GVK.Kind {
				found = true
				continue
			}

			if !found {
				continue
			}

			manifestPath = append(manifestPath, p)
		}

		cachedPath := strings.Join(manifestPath, ".")
		if _, ok := cache[cachedPath]; ok {
			continue
		}

		cache[cachedPath] = true

		v, err := d.Properties.Value(manifestPath)
		if err != nil {
			return nil, errors.Wrapf(err, "retrieve values for %s", strings.Join(manifestPath, "."))
		}

		dv := documentValues{
			setter: item.Name,
			value:  v,
		}

		p := strings.Join(item.Path, ".")
		m[p] = dv
	}

	return m, nil
}

func (d *Document) resolvedPaths() (*Locals, error) {
	nn := NewNode("root", d.root)

	locals := NewLocals(d.GVK.Kind)

	paths := d.Properties.Paths(d.GVK)
	for _, path := range paths {
		sr, err := nn.Search(path.Path...)
		if err != nil {
			return nil, errors.Wrapf(err, "search path %s", strings.Join(path.Path, "."))
		}

		manifestPath := sr.MatchedPath[4:]
		if path.Path[len(path.Path)-1] != manifestPath[len(manifestPath)-1] {
			continue
		}

		// spew.Dump(path, sr)

		var paramName bytes.Buffer
		for i := range manifestPath {
			part := manifestPath[i]
			if i > 0 {
				part = strings.Title(part)
			}

			paramName.WriteString(part)
		}

		_, err = d.Properties.Value(manifestPath)
		if err != nil {
			return nil, errors.Wrapf(err, "retrieve manifest values for %s", strings.Join(path.Path, "."))
		}

		// if err := comp.AddParam(paramName.String(), v); err != nil {
		// 	return nil, errors.Wrapf(err, "add param %s to component", paramName.String())
		// }

		k := strings.Join(sr.MatchedPath[:len(sr.MatchedPath)-1], ".")
		// fmt.Println(k, "=", v)
		entry := LocalEntry{
			Path:      k,
			Setter:    sr.Setter,
			ParamName: paramName.String(),
		}

		locals.Add(entry)
	}

	// return locals, nil
	return nil, nil
}
