package main

import (
	"flag"

	"github.com/bryanl/woowoo/yaml2jsonnet"
	"github.com/sirupsen/logrus"
)

var (
	libRoot = "/tmp/k8s.libsonnet"
)

func main() {
	var source string
	flag.StringVar(&source, "source", "", "Kubernetes manifest")
	flag.Parse()

	conversion, err := yaml2jsonnet.NewConversion(source, libRoot)
	if err != nil {
		logrus.WithError(err).Fatal("initialize conversion")
	}

	if err := conversion.Process(); err != nil {
		logrus.WithError(err).Fatal("process document")
	}
}
