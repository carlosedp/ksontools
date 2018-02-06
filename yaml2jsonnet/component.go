package yaml2jsonnet

import (
	"bytes"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/pkg/errors"
)

var (
	ksonnetImport = Declaration{Name: "k", Value: NewDeclarationImport("k.libsonnet")}
)

type Declaration struct {
	Name  string
	Value DeclarationValue
}

type DeclarationValue interface {
	Node() ast.Node
}

type DeclarationImport struct {
	s string
}

func NewDeclarationImport(s string) *DeclarationImport {
	return &DeclarationImport{s: s}
}

func (di *DeclarationImport) Node() ast.Node {
	n := &ast.Import{
		File: &ast.LiteralString{
			Kind:  ast.StringDouble,
			Value: di.s,
		},
	}

	return n
}

type DeclarationString struct {
	s string
}

func NewDeclarationString(s string) *DeclarationString {
	return &DeclarationString{s: s}
}

func (ds *DeclarationString) Node() ast.Node {
	n := &ast.LiteralString{
		Kind:  ast.StringDouble,
		Value: ds.s,
	}

	return n
}

type DeclarationApply struct {
	s string
}

func NewDeclarationApply(s string) *DeclarationApply {
	return &DeclarationApply{s: s}
}

func (da *DeclarationApply) Node() ast.Node {
	call := nodemaker.NewCall(da.s)
	return call.Node()
}

type DeclarationNoder struct {
	noder nodemaker.Noder
}

func NewDeclarationNoder(noder nodemaker.Noder) *DeclarationNoder {
	return &DeclarationNoder{noder: noder}
}

func (dn *DeclarationNoder) Node() ast.Node {
	return dn.noder.Node()
}

// Component generates a component with parameters.
type Component struct {
	declarations []Declaration
	params       *nodemaker.Object
}

func NewComponent() *Component {
	return &Component{
		declarations: make([]Declaration, 0),
		params:       nodemaker.NewObject(),
	}
}

func (c *Component) AddDeclaration(d Declaration) {
	c.declarations = append(c.declarations, d)
}

// Declarations prints declarations for a component.
func (c *Component) Declarations(obj ast.Node) (*ast.Local, error) {
	decs := make([]*ast.Local, len(c.declarations))

	for i, d := range c.declarations {
		decs[i] = genDeclaration(d)
	}

	var root *ast.Local
	var prev *ast.Local

	for i, dec := range decs {
		if i == 0 {
			root = dec
			prev = dec
			continue
		}
		prev.Body = dec
		prev = dec
	}

	prev.Body = obj

	return root, nil
}

func genDeclaration(decl Declaration) *ast.Local {
	id := ast.Identifier(decl.Name)

	return &ast.Local{
		Binds: ast.LocalBinds{
			ast.LocalBind{
				Variable: id,
				Body:     decl.Value.Node(),
			},
		},
	}
}

// Generate generates a component.
func (c *Component) Generate(obj ast.Node) (string, error) {
	// write out header
	header := genDeclaration(ksonnetImport)

	paramsDecl := Declaration{
		Name:  "params",
		Value: NewDeclarationNoder(c.params),
	}

	params := genDeclaration(paramsDecl)
	header.Body = params

	decls, err := c.Declarations(obj)
	if err != nil {
		return "", errors.Wrap(err, "create libsonnet declarations")
	}

	params.Body = decls

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, header); err != nil {
		return "", errors.Wrap(err, "create jsonnet")
	}

	return buf.String(), nil
}

func (c *Component) AddParam(name string, value interface{}) error {
	return addParam(c.params, name, value)
}

func addParam(parent *nodemaker.Object, name string, value interface{}) error {
	keyName := nodemaker.InheritedKey(name)

	switch t := value.(type) {
	default:
		return errors.Errorf("unable to handle param %s of type %T", name, t)
	case string:
		parent.Set(keyName, nodemaker.NewStringDouble(t))
	case int:
		parent.Set(keyName, nodemaker.NewInt(t))
	case float64:
		parent.Set(keyName, nodemaker.NewFloat(t))
	case []interface{}:
		var nodes []nodemaker.Noder
		for _, elem := range t {
			node, err := createNode(elem)
			if err != nil {
				return errors.Wrap(err, "add to array")
			}
			nodes = append(nodes, node)
		}

		parent.Set(keyName, nodemaker.NewArray(nodes))
	case map[interface{}]interface{}:
		o := nodemaker.NewObject()
		for k := range t {
			n, ok := k.(string)
			if !ok {
				return errors.Errorf("object key is not string (%T)", k)
			}

			if err := addParam(o, n, t[k]); err != nil {
				return errors.Wrap(err, "add child key")
			}
		}

		parent.Set(keyName, o)
	}

	return nil
}

func createNode(item interface{}) (nodemaker.Noder, error) {
	switch t := item.(type) {
	default:
		return nil, errors.Errorf("unable to create node of type %T", t)
	case int:
		return nodemaker.NewInt(t), nil
	case string:
		return nodemaker.NewStringDouble(t), nil
	case float64:
		return nodemaker.NewFloat(t), nil
	case map[interface{}]interface{}:
		parent := nodemaker.NewObject()
		for k := range t {
			name, ok := k.(string)
			if !ok {
				return nil, errors.Errorf("object key is not string (%T)", k)
			}
			if err := addParam(parent, name, t[k]); err != nil {
				return nil, errors.Wrap(err, "add child key")
			}
		}

		return parent, nil
	}
}
