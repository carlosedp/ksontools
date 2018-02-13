package main

import (
	"github.com/bryanl/woowoo/docgen"
	"github.com/sirupsen/logrus"
)

func main() {
	path := "/Users/bryan/Development/heptio/ksonnet/ksonnet-playground/new-ctors/k8s.libsonnet"
	outPath := "/Users/bryan/go/src/github.com/bryanl/woowoo/k8sdocs"

	if err := docgen.Generate(path, outPath); err != nil {
		logrus.WithError(err).Fatal("create ksonnet lib docs")
	}
}
