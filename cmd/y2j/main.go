package main

import (
	"flag"
	"io/ioutil"

	"github.com/bryanl/woowoo/yaml2jsonnet"
	"github.com/sirupsen/logrus"
)

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "version", true, "Verbose mode")

	var k8slib string
	flag.StringVar(&k8slib, "k8slib", "k8s.libsonnet", "Path to k8s.libsonnet")

	flag.Parse()

	if !verbose {
		logrus.SetOutput(ioutil.Discard)
	}

	if flag.NArg() != 1 {
		logrus.Fatal("must supply source")
	}

	source := flag.Arg(0)

	conversion, err := yaml2jsonnet.NewConversion(source, k8slib)
	if err != nil {
		logrus.WithError(err).Fatal("initialize conversion")
	}

	if err := conversion.Process(); err != nil {
		logrus.WithError(err).Fatal("process document")
	}
}
