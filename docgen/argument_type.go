package docgen

import (
	"github.com/gobuffalo/plush"

	rice "github.com/GeertJohan/go.rice"
)

type argumentType struct {
	propertyName string
	typeName     string
	typeDef      string
}

func (at *argumentType) ToDoc() (string, error) {
	box, err := rice.FindBox("templates")
	if err != nil {
		return "", err
	}

	tmpl, err := box.String("argument_type.html")
	if err != nil {
		return "", err
	}

	ctx := plush.NewContext()
	ctx.Set("propertyName", at.propertyName)
	ctx.Set("typeName", at.typeName)
	ctx.Set("typeDef", at.typeDef)

	return plush.Render(tmpl, ctx)
}
