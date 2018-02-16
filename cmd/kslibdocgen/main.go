package main

import (
	"flag"
	"strings"

	"github.com/bryanl/woowoo/docgen"
	"github.com/sirupsen/logrus"
)

func main() {
	var groups arrayFlags
	flag.Var(&groups, "groups", "Groups to render. If blank, it will render all groups")

	var path string
	flag.StringVar(&path, "path", "/Users/bryan/Development/heptio/ksonnet/ksonnet-playground/new-ctors/k8s.libsonnet", "Path to ksonnet")

	var outPath string
	flag.StringVar(&outPath, "outPath", "/Users/bryan/go/src/github.com/bryanl/woowoo/k8sdocs", "Output path")

	flag.Parse()

	if err := docgen.Generate(path, outPath, []string(groups)...); err != nil {
		logrus.WithError(err).Fatal("create ksonnet lib docs")
	}
}

type arrayFlags []string

func (f *arrayFlags) String() string {
	return strings.Join(*f, ",")
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}
