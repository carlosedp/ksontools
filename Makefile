.PHONY = install-kscomp

install-kscomp:
	mkdir -p ${HOME}/.config/ksonnet/plugins/kscomp
	cp hack/kscomp-plugin.yaml ${HOME}/.config/ksonnet/plugins/kscomp/plugin.yaml
	go build -o ${HOME}/.config/ksonnet/plugins/kscomp/kscomp github.com/bryanl/woowoo/cmd/kscomp