package yaml2jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
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
	buildConstructors map[string][]ctorArgument
	componentName     string
}

// NewDocument creates an instance of Document.
func NewDocument(componentName string, r io.Reader, root ast.Node) (*Document, error) {
	obj, ok := root.(*astext.Object)
	if !ok {
		return nil, errors.New("root is not an *ast.Object")
	}

	doc := &Document{
		root:          obj,
		componentName: componentName,
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

	resolvedPaths, err := doc.generateResolvedPaths()
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

type localBlock struct {
	locals []*nm.Local
}

func newLocalBlock() *localBlock {
	return &localBlock{
		locals: make([]*nm.Local, 0),
	}
}

func (lb *localBlock) add(local *nm.Local) {
	lb.locals = append(lb.locals, local)
}

func (lb *localBlock) node(body nm.Noder) nm.Noder {
	for i, local := range lb.locals {
		if i == len(lb.locals)-1 {
			local.Body = body
			continue
		}

		local.Body = lb.locals[i+1]
	}

	return lb.locals[0]
}

// GenerateComponent generates a component
func (d *Document) GenerateComponent() (string, error) {
	componentName := d.componentName
	logrus.WithField("componentName", componentName).Info("generating component")

	lb := newLocalBlock()
	lb.add(d.importParams())
	lb.add(createLocal("k", nm.NewImport("k.libsonnet")))

	mixins := d.buildMixins()
	for _, mixin := range mixins {
		lb.add(mixin)
	}

	objectCtorName, objectCtorFn := d.buildObjectCtor()
	lb.add(createLocal(objectCtorName, objectCtorFn))
	lb.add(createLocal(componentName, d.buildObject()))

	body := nm.NewVar(componentName)
	node := lb.node(body)

	return d.render(node.Node())
}

// ParamsUpdater is an interface for updating params.
type ParamsUpdater interface {
	Update(componentName string, params map[string]interface{}) error
}

// UpdateParams updates params.
func (d *Document) UpdateParams(pu ParamsUpdater) error {
	m := make(map[string]interface{})
	for _, ctorArugments := range d.buildConstructors {
		for _, ca := range ctorArugments {
			m[ca.paramName] = ca.paramValue
		}
	}

	return pu.Update(d.componentName, m)
}

func (d *Document) genParams() map[string]interface{} {
	m := make(map[string]interface{})
	for ns, dv := range d.resolvedPaths {
		name := paramName(ns)
		m[name] = dv.value
	}

	return m
}

func createLocal(name string, value nm.Noder) *nm.Local {
	return nm.NewLocal(name, value, nil)
}

func (d *Document) importParams() *nm.Local {
	cc := nm.NewCallChain(
		nm.NewVar("std"),
		nm.NewApply(nm.NewIndex("extVar"), []nm.Noder{
			nm.NewStringDouble("__ksonnet/params"),
		}, nil),
		nm.NewIndex("components"),
		nm.NewIndex(d.componentName),
	)

	return createLocal("params", cc)
}

func (d *Document) buildObject() nm.Noder {
	componentName := d.componentName
	objectCtorName := genObjectCtorName(componentName)
	apply := nm.NewApply(nm.NewVar(objectCtorName), []nm.Noder{nm.NewVar("params")}, nil)
	return apply
}

func genObjectCtorName(componentName string) string {
	return strcase.ToLowerCamel(fmt.Sprintf("%s_%s", "create", componentName))
}

func (d *Document) buildObjectCtor() (string, nm.Noder) {
	componentName := d.componentName
	objectCtorName := genObjectCtorName(componentName)

	pathPrefix := append([]string{"k"}, d.GVK.Path()...)
	objectCtorPath := strings.Join(append(pathPrefix, "new"), ".")
	objectCtorCall := nm.ApplyCall(objectCtorPath)

	nodes := []nm.Noder{objectCtorCall}

	locals := newLocalBlock()

	for _, ns := range d.paths() {
		ctorArguments := d.buildConstructors[ns]
		var args []nm.Noder
		for _, ca := range ctorArguments {
			args = append(args, nm.NewCall(fmt.Sprintf("params.%s", ca.paramName)))
		}

		objectName := mixinObjectName(ns)
		ctorName := mixinConstructorName(ns)
		ctorApply := nm.ApplyCall(ctorName, args...)

		local := createLocal(objectName, ctorApply)
		locals.add(local)
		nodes = append(nodes, nm.NewVar(objectName))
	}

	combiner := nm.Combine(nodes...)
	node := locals.node(combiner)

	objectCtorFn := nm.NewFunction([]string{"params"}, node)

	return objectCtorName, objectCtorFn
}

func (d *Document) buildMixins() []*nm.Local {
	var locals []*nm.Local

	var names []string
	for ns := range d.buildConstructors {
		names = append(names, ns)
	}
	sort.Strings(names)

	for _, ns := range names {
		ctorArguments := d.buildConstructors[ns]
		fnName := mixinConstructorName(ns)

		links := []nm.Chainable{
			nm.NewVar("k"),
			nm.NewCall(ns),
		}

		var args = []string{}
		for _, ca := range ctorArguments {
			arg := strings.TrimPrefix(ca.setter, "with")
			arg = strings.ToLower(arg)
			args = append(args, arg)

			links = append(
				links,
				nm.NewApply(
					nm.NewIndex(ca.setter),
					[]nm.Noder{nm.NewVar(arg)},
					nil))
		}

		fn := nm.NewFunction(args, nm.NewCallChain(links...))
		locals = append(locals, createLocal(fnName, fn))
	}

	return locals
}

func paramName(s string) string {
	nameParts := mixinNameParts(s)

	var name string

	for i, r := range nameParts[0] {
		if i == 0 {
			name += string(r)
			continue
		}

		if unicode.IsUpper(r) {
			name += string(unicode.ToLower(r))
		}
	}

	for _, s := range nameParts[1:] {
		name += strings.Title(s)
	}

	return name
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

func (d *Document) paths() []string {
	var names []string
	for ns := range d.buildConstructors {
		names = append(names, ns)
	}
	sort.Strings(names)

	return names
}

func (d *Document) generateResolvedPaths() (map[string]documentValues, error) {
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
