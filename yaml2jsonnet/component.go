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
}

func NewComponent() *Component {
	return &Component{
		declarations: make([]Declaration, 0),
	}
}

func (c *Component) AddDeclaration(d Declaration) {
	c.declarations = append(c.declarations, d)
}

// Declarations prints declarations for a component.
func (c *Component) Declarations(node ast.Node) (string, error) {
	curDecs := append([]Declaration{ksonnetImport}, c.declarations...)

	decs := make([]*ast.Local, len(curDecs))

	for i, d := range curDecs {
		id := ast.Identifier(d.Name)

		decs[i] = &ast.Local{
			Binds: ast.LocalBinds{
				ast.LocalBind{
					Variable: id,
					Body:     d.Value.Node(),
				},
			},
		}
	}

	var root *ast.Local
	var prev *ast.Local

	for i := range decs {
		dec := decs[i]
		if i == 0 {
			root = dec
			prev = dec
			continue
		}

		prev.Body = dec
		prev = dec
	}

	prev.Body = node

	var buf bytes.Buffer

	if err := printer.Fprint(&buf, root); err != nil {
		return "", errors.Wrap(err, "create jsonnet")
	}

	return buf.String(), nil
}

// Generate generates a component.
func (c *Component) Generate(obj ast.Node) (string, error) {
	s, err := c.Declarations(obj)
	if err != nil {
		return "", errors.Wrap(err, "create libsonnet declarations")
	}

	return s, nil
}
