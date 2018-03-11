package env

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksutil"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Apply applies components to a cluster.
// func Apply(ksApp ksutil.SuperApp, envName string, components []string) error {
// 	objects, err := buildObjects(ksApp, envName, components)
// 	if err != nil {
// 		return err
// 	}

// 	// TODO: this is hackish. create better semantics around apply
// 	c := kubecfg.ApplyCmd{}

// 	return nil
// }

// Show shows YAML rendered for an environment.
func Show(ksApp ksutil.SuperApp, envName string, components []string) error {
	objects, err := buildObjects(ksApp, envName, components)
	if err != nil {
		return err
	}

	if err := ksutil.Fprint(os.Stdout, objects, "yaml"); err != nil {
		return errors.Wrap(err, "print YAML")
	}

	return nil
}

func buildObjects(ksApp ksutil.SuperApp, envName string, components []string) ([]*unstructured.Unstructured, error) {
	namespaces, err := component.NamespacesFromEnv(ksApp, envName)
	if err != nil {
		return nil, errors.Wrap(err, "find namespaces")
	}

	var objects []*unstructured.Unstructured
	for _, ns := range namespaces {
		paramsStr, err := buildEnvParam(ksApp, envName, ns)
		if err != nil {
			return nil, err
		}

		members, err := ns.Components()
		if err != nil {
			return nil, errors.Wrap(err, "find components")
		}

		for _, c := range members {
			if !isAvailable(c.Name(), components) {
				continue
			}

			o, err := c.Objects(paramsStr)
			if err != nil {
				return nil, errors.Wrap(err, "get objects")
			}
			objects = append(objects, o...)
		}
	}

	return objects, nil
}

func isAvailable(name string, components []string) bool {
	if len(components) == 0 {
		return true
	}

	return stringInSlice(name, components)
}

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}

func buildEnvParam(ksApp ksutil.SuperApp, envName string, ns component.Namespace) (string, error) {
	paramsStr, err := ns.ResolvedParams()
	if err != nil {
		return "", err
	}

	envParamsPath := filepath.Join(ksApp.Root(), "environments", envName, "params.libsonnet")
	b, err := afero.ReadFile(ksApp.Fs(), envParamsPath)
	if err != nil {
		return "", err
	}

	envParams := upgradeParams(string(b))

	vm := jsonnet.MakeVM()
	vm.ExtCode("__ksonnet/params", paramsStr)
	return vm.EvaluateSnippet("snippet", string(envParams))
}

var (
	reParamSwap = regexp.MustCompile(`(?m)import "\.\.\/\.\.\/components\/params\.libsonnet"`)
)

// upgradeParams replaces relative params imports with an extVar to handle
// multiple component namespaces.
// NOTE: It warns when it makes a change. This serves as a temporary fix until
// ksonnet generates the correct file.
func upgradeParams(in string) string {
	return reParamSwap.ReplaceAllLiteralString(in, `std.extVar("__ksonnet/params")`)
}
