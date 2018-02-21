package yaml2jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/sirupsen/logrus"

	"github.com/go-yaml/yaml"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

// Document creates a ksonnet document for describing a resource.
type Document struct {
	Properties        Properties
	GVK               GVK
	root              *astext.Object
	resolvedPaths     map[string]documentValues
	buildConstructors map[string][]string
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

	resolvedPaths, err := doc.resolvedPaths2()
	if err != nil {
		return nil, errors.Wrap(err, "resolve document paths")
	}

	doc.resolvedPaths = resolvedPaths

	ctors, err := buildConstructors(resolvedPaths)
	if err != nil {
		return nil, errors.Wrap(err, "build object constructors")
	}

	doc.buildConstructors = ctors

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

type localBlock struct {
	locals []*nodemaker.Local
}

func newLocalBlock() *localBlock {
	return &localBlock{
		locals: make([]*nodemaker.Local, 0),
	}
}

func (lb *localBlock) add(local *nodemaker.Local) {
	lb.locals = append(lb.locals, local)
}

func (lb *localBlock) node(body nodemaker.Noder) nodemaker.Noder {
	for i, local := range lb.locals {
		if i == len(lb.locals)-1 {
			local.Body = body
			continue
		}

		local.Body = lb.locals[i+1]
	}

	return lb.locals[0]
}

// GenerateComponent2 generates a component
func (d *Document) GenerateComponent2(componentName string) (string, error) {
	logrus.WithField("componentName", componentName).Info("generating component")

	lb := newLocalBlock()
	lb.add(importParams(componentName))
	lb.add(createLocal("k", nodemaker.NewImport("k.libsonnet")))

	mixins := d.buildMixins()
	for _, mixin := range mixins {
		lb.add(mixin)
	}

	objectCtorName, objectCtorFn := d.buildObject(componentName)
	lb.add(createLocal(objectCtorName, objectCtorFn))

	body := nodemaker.NewObject()

	node := lb.node(body)
	return d.render(node.Node())
}

func createLocal(name string, value nodemaker.Noder) *nodemaker.Local {
	return nodemaker.NewLocal(name, value, nil)
}

func importParams(componentName string) *nodemaker.Local {
	cc := nodemaker.NewCallChain(
		nodemaker.NewVar("std"),
		nodemaker.NewApply(nodemaker.NewIndex("extVar"), []nodemaker.Noder{
			nodemaker.NewStringDouble("__ksonnet/params"),
		}, nil),
		nodemaker.NewIndex("components"),
		nodemaker.NewIndex(componentName),
	)

	return createLocal("params", cc)
}

func (d *Document) buildObject(componentName string) (string, nodemaker.Noder) {
	objectCtorName := strcase.ToLowerCamel(fmt.Sprintf("%s_%s", "create", componentName))

	pathPrefix := append([]string{"k"}, d.GVK.Path()...)
	objectCtorPath := strings.Join(append(pathPrefix, "new"), ".")
	objectCtorCall := nodemaker.ApplyCall(objectCtorPath)

	nodes := []nodemaker.Noder{objectCtorCall}

	locals := newLocalBlock()
	for ns := range d.buildConstructors {
		objectName := mixinObjectName(ns)
		ctorName := mixinConstructorName(ns)
		ctorApply := nodemaker.ApplyCall(ctorName)

		local := createLocal(objectName, ctorApply)
		locals.add(local)
		nodes = append(nodes, nodemaker.NewVar(objectName))
	}

	combiner := nodemaker.Combine(nodes...)
	node := locals.node(combiner)

	objectCtorFn := nodemaker.NewFunction([]string{"params"}, node)

	return objectCtorName, objectCtorFn
}

func (d *Document) buildMixins() []*nodemaker.Local {
	var locals []*nodemaker.Local
	for ns, setters := range d.buildConstructors {
		fnName := mixinConstructorName(ns)

		links := []nodemaker.Chainable{
			nodemaker.NewVar("k"),
			nodemaker.NewCall(ns),
		}

		var args = []string{}
		for _, setter := range setters {
			arg := strings.TrimPrefix(setter, "with")
			arg = strings.ToLower(arg)
			args = append(args, arg)

			links = append(
				links,
				nodemaker.NewApply(
					nodemaker.NewIndex(setter),
					[]nodemaker.Noder{nodemaker.NewVar(arg)},
					nil))
		}

		fn := nodemaker.NewFunction(args, nodemaker.NewCallChain(links...))
		locals = append(locals, createLocal(fnName, fn))
	}

	return locals
}

func mixinConstructorName(ns string) string {
	nameParts := mixinNameParts(ns)
	fnName := strcase.ToLowerCamel(fmt.Sprintf("create_%s", strings.Join(nameParts, "_")))
	return fnName
}

func mixinObjectName(ns string) string {
	nameParts := mixinNameParts(ns)
	objectName := strcase.ToLowerCamel(strings.Join(nameParts, "_"))
	return objectName
}

func mixinNameParts(ns string) []string {
	parts := strings.Split(ns, ".")
	mixinIndex := -1
	for i, part := range parts {
		if part == "mixin" {
			mixinIndex = i
		}
	}

	nameParts := make([]string, len(parts))
	copy(nameParts, parts)
	if mixinIndex >= 0 {
		nameParts = append(nameParts[:mixinIndex], nameParts[mixinIndex+1:]...)
	}

	return nameParts[2:]
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
