package main

import (
	"flag"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	jsonnetutil "github.com/ksonnet/ksonnet/pkg/util/jsonnet"
	"github.com/sirupsen/logrus"
)

func main() {
	var source string
	flag.StringVar(&source, "source", "", "jsonnet source")
	flag.Parse()

	node, err := jsonnetutil.Import(source)
	if err != nil {
		logrus.WithError(err).Fatal("import jsonnet")
	}

	if err := printer.Fprint(os.Stdout, node); err != nil {
		logrus.WithError(err).Fatal("print document")
	}

	spew.Dump(node)
}
