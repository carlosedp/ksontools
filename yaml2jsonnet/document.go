package yaml2jsonnet

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"

	"github.com/go-yaml/yaml"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Document struct {
	Properties map[string]interface{}
	GVK        GVK
	root       ast.Node
}

func NewDocument(r io.Reader, root ast.Node) (*Document, error) {
	doc := &Document{
		Properties: make(map[string]interface{}),
		root:       root,
	}

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
			doc.Properties[k] = v
		}
	}

	gvk, err := ts.GVK()
	if err != nil {
		return nil, errors.Wrap(err, "type spec is invalid")
	}

	doc.GVK = gvk

	return doc, nil
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

	obj, err := FindType(d.GVK, d.root)
	if err != nil {
		return "", errors.Wrap(err, "find root node")
	}

	var mixinNames []string

	root := NewNode(d.GVK.Kind, obj)

	var names []string
	for k := range d.Properties {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, name := range names {
		value := d.Properties[name]
		logger := logrus.WithFields(logrus.Fields{
			"name": name,
		})

		var builders []string

		switch t := value.(type) {
		default:
			logger.WithField("type", fmt.Sprintf("%T", t)).
				Warn("not sure what to do with this")
		case map[interface{}]interface{}:
			node, err := root.Property(name)
			if err != nil {
				return "", errors.Wrapf(err, "inspect property %s", name)
			}

			for k, v := range t {
				k1 := k.(string)
				setter, err := node.FindFunction(name, k1)
				if err != nil {
					logger.Warnf("%s is a mixin", k1)
					continue
				}

				comp.AddParam(k1, v)

				logger.Infof("setter for %s is %s", k1, setter)

				builders = append(builders, fmt.Sprintf("%s(%s)", setter, k1))
			}

			if node.IsMixin && len(builders) > 0 {
				method := fmt.Sprintf("%s.mixin.%s.%s",
					d.GVK.Kind,
					node.name,
					strings.Join(builders, "."))

				val := NewDeclarationApply(method)

				mixinName := fmt.Sprintf("%s%s", d.GVK.Kind, strings.Title(node.name))

				decl := Declaration{
					Name:  mixinName,
					Value: val,
				}
				comp.AddDeclaration(decl)

				mixinNames = append(mixinNames, mixinName)
			}
		}
	}

	nodeInit := fmt.Sprintf("init%s", strings.Title(d.GVK.Kind))

	comp.AddDeclaration(Declaration{
		Name:  nodeInit,
		Value: NewDeclarationNoder(nodemaker.ApplyCall(fmt.Sprintf("%s.new", d.GVK.Kind))),
	})

	n := nodemaker.NewVar(nodeInit)

	var left nodemaker.Noder = n
	for _, name := range mixinNames {
		left = nodemaker.NewBinary(left, nodemaker.NewVar(name), nodemaker.BopPlus)
	}

	s, err := comp.Generate(left.Node())
	if err != nil {
		return "", errors.Wrap(err, "generate component")
	}

	return s, nil
}
