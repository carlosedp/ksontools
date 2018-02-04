package yaml2jsonnet

import (
	"bytes"
	"fmt"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	keyName := nodemaker.InheritedKey(name)

	switch t := value.(type) {
	default:
		logrus.WithFields(logrus.Fields{
			"name":  name,
			"value": fmt.Sprintf("%#v", t),
			"type":  fmt.Sprintf("%T", t),
		}).Error("adding param")
	case string:
		c.params.Set(keyName, nodemaker.NewStringDouble(t))
	case int:
		c.params.Set(keyName, nodemaker.NewInt(t))
	case float64:
		c.params.Set(keyName, nodemaker.NewFloat(t))
	case map[interface{}]interface{}:
		o := nodemaker.NewObject()
		for k, v := range t {
			k1, ok := k.(string)
			if !ok {
				return errors.Errorf("param %s is not a string key", name)
			}

			nestedKeyName := nodemaker.InheritedKey(k1)
			switch t1 := v.(type) {
			default:
				logrus.WithFields(logrus.Fields{
					"name":  k1,
					"value": fmt.Sprintf("%#v", t1),
					"type":  fmt.Sprintf("%T", t1),
				}).Error("adding object field param")
			case string:
				o.Set(nestedKeyName, nodemaker.NewStringDouble(t1))
			case int:
				o.Set(nestedKeyName, nodemaker.NewInt(t1))
			case float64:
				o.Set(nestedKeyName, nodemaker.NewFloat(t1))
			}
		}

		c.params.Set(keyName, o)
	}
	return nil
}
