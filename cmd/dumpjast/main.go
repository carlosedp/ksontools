package main

import (
	"flag"
	"os"

	"github.com/bryanl/woowoo/yaml2jsonnet"
	"github.com/davecgh/go-spew/spew"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/sirupsen/logrus"
)

func main() {
	var source string
	flag.StringVar(&source, "source", "", "jsonnet source")
	flag.Parse()

	node, err := yaml2jsonnet.ImportJsonnet(source)
	if err != nil {
		logrus.WithError(err).Fatal("import jsonnet")
	}

	if err := printer.Fprint(os.Stdout, node); err != nil {
		logrus.WithError(err).Fatal("print document")
	}

	spew.Dump(node)
}
