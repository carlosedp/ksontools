package docgen

import (
	"github.com/bryanl/woowoo/yaml2jsonnet"

	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

// Import imports jsonnet
func Import(k8sLibPath, docsPath string) error {
	node, err := yaml2jsonnet.ImportJsonnet(k8sLibPath)
	if err != nil {
		return errors.Wrap(err, "parse and evaluate source")
	}

	dg, err := New(node, docsPath)
	if err != nil {
		return errors.Wrap(err, "create Docgen")
	}

	return dg.Generate()
}

type groupKind struct {
	group string
	kind  string
}

// Docgen is a documentation generator for k8s.
type Docgen struct {
	node ast.Node
	hugo *hugo

	versionLookup map[groupKind][]string
}

// New creates an instance of Docgen.
func New(node ast.Node, docsPath string) (*Docgen, error) {
	h, err := newHugo(docsPath)
	if err != nil {
		return nil, err
	}

	return &Docgen{
		node:          node,
		hugo:          h,
		versionLookup: make(map[groupKind][]string),
	}, nil
}

// Generate generates documentation.
func (dg *Docgen) Generate() error {
	err := iterateObject(dg.node, dg.generateGroup)
	return errors.Wrap(err, "iterate over groups")
}

func (dg *Docgen) generateGroup(name string, node ast.Node) error {
	fm := newGroupFrontMatter(name)

	if err := dg.hugo.writeGroup(name, fm); err != nil {
		return errors.Wrap(err, "write group")
	}

	fn := func(version string, node ast.Node) error {
		return dg.generateVersion(name, version, node)
	}

	err := iterateObject(node, fn)
	if err != nil {
		return errors.Wrapf(err, "iterate over group %s", name)
	}

	for gk, versions := range dg.versionLookup {
		fm := newKindFrontMatter(gk.group, gk.kind, versions)

		if err := dg.hugo.writeKind(gk.group, gk.kind, fm); err != nil {
			return errors.Wrapf(err, "write kind %s/%s", gk.group, gk.kind)
		}

		for _, version := range versions {
			if err := dg.hugo.writeVersionedKind(gk.group, version, gk.kind); err != nil {
				return errors.Wrapf(err, "write versioned kind %s/%s/%s", gk.group, version, gk.kind)
			}
		}
	}

	return nil
}

func (dg *Docgen) generateVersion(group, version string, node ast.Node) error {
	fn := func(name string, node ast.Node) error {
		return dg.generateKind(group, version, name, node)
	}

	if err := iterateObject(node, fn); err != nil {
		return err
	}

	return nil
}

func (dg *Docgen) generateKind(group, version, kind string, node ast.Node) error {
	if kind == "apiVersion" {
		return nil
	}

	gk := groupKind{group: group, kind: kind}

	_, ok := dg.versionLookup[gk]
	if !ok {
		dg.versionLookup[gk] = make([]string, 0)
	}

	dg.versionLookup[gk] = append(dg.versionLookup[gk], version)

	obj, ok := node.(*ast.Object)
	if !ok {
		return errors.Errorf("%s/%s/%s is not an object", group, version, kind)
	}

	for _, of := range obj.Fields {
		id := string(*of.Id)
		if of.Method == nil {
			continue
		}

		if err := dg.hugo.writeProperty(group, version, kind, id); err != nil {
			return err
		}
	}

	return nil
}

func iterateObject(node ast.Node, fn func(string, ast.Node) error) error {
	if node == nil {
		return errors.New("node was nil")
	}

	obj, ok := node.(*ast.Object)
	if !ok {
		return errors.New("node was not an object")
	}

	for _, of := range obj.Fields {
		if of.Hide == ast.ObjectFieldInherit {
			continue
		}

		if of.Kind == ast.ObjectLocal {
			continue
		}

		id := string(*of.Id)
		if id == "hidden" {
			continue
		}

		if err := fn(id, of.Expr2); err != nil {
			return err
		}
	}

	return nil
}
