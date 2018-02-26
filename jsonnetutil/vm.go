package jsonnetutil

import (
	"fmt"
	"net/http"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/ksonnet/ksonnet/utils"
)

// VMFactory produces jsonnet VM.
type VMFactory struct {
	JPath    []string
	ExtVars  map[string]string
	TLAVars  map[string]string
	ExtCodes map[string]string

	Resolver   string
	FailAction string
}

// VM constructs a new jsonnet.VM.
func (vf *VMFactory) VM() (*jsonnet.VM, error) {
	vm := jsonnet.MakeVM()
	importer := jsonnet.FileImporter{
		JPaths: []string{},
	}

	for _, p := range vf.JPath {
		log.Debugln("Adding jsonnet search path", p)
		importer.JPaths = append(importer.JPaths, p)
	}

	vm.Importer(&importer)

	for k, v := range vf.ExtVars {
		vm.ExtVar(k, v)
	}

	for k, v := range vf.TLAVars {
		vm.TLAVar(k, v)
	}

	for k, v := range vf.ExtCodes {
		vm.ExtCode(k, v)
	}

	resolver, err := vf.buildResolver()
	if err != nil {
		return nil, err
	}
	RegisterNativeFuncs(vm, resolver)

	return vm, nil
}

func (vf *VMFactory) buildResolver() (Resolver, error) {
	ret := resolverErrorWrapper{}

	switch vf.FailAction {
	case "ignore":
		ret.OnErr = func(error) error { return nil }
	case "warn":
		ret.OnErr = func(err error) error {
			log.Warning(err.Error())
			return nil
		}
	case "error":
		ret.OnErr = func(err error) error { return err }
	default:
		return nil, fmt.Errorf("Unknown resolve failure type: %s", vf.FailAction)
	}

	switch vf.Resolver {
	case "noop":
		ret.Inner = NewIdentityResolver()
	case "registry":
		ret.Inner = NewRegistryResolver(&http.Client{
			Transport: utils.NewAuthTransport(http.DefaultTransport),
		})
	default:
		return nil, fmt.Errorf("Unknown resolver type: %s", vf.Resolver)
	}

	return &ret, nil
}

type resolverErrorWrapper struct {
	Inner Resolver
	OnErr func(error) error
}

func (r *resolverErrorWrapper) Resolve(image *ImageName) error {
	err := r.Inner.Resolve(image)
	if err != nil {
		err = r.OnErr(err)
	}
	return err
}
