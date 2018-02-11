package main

import (
	"flag"
	"io/ioutil"

	"github.com/bryanl/woowoo/yaml2jsonnet"
	"github.com/sirupsen/logrus"
)

var (
	libRoot = "/tmp/k8s.libsonnet"
)

func main() {

	var source string
	flag.StringVar(&source, "source", "", "Kubernetes manifest")

	var verbose bool
	flag.BoolVar(&verbose, "version", true, "Verbose mode")
	flag.Parse()

	if !verbose {
		logrus.SetOutput(ioutil.Discard)
	}

	conversion, err := yaml2jsonnet.NewConversion(source, libRoot)
	if err != nil {
		logrus.WithError(err).Fatal("initialize conversion")
	}

	if err := conversion.Process(); err != nil {
		logrus.WithError(err).Fatal("process document")
	}
}
