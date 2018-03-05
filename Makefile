.PHONY = install-kscomp update-docgen doc-groups doc

install-kscomp:
	mkdir -p ${HOME}/.config/ksonnet/plugins/kscomp
	cp hack/kscomp-plugin.yaml ${HOME}/.config/ksonnet/plugins/kscomp/plugin.yaml
	go build -o ${HOME}/.config/ksonnet/plugins/kscomp/kscomp github.com/bryanl/woowoo/cmd/kscomp

update-docgen:
	cd docgen && \
	rice embed-go -v

doc-groups:
	go run cmd/kslibdocgen/main.go --path tmp/k8s.libsonnet --groups apps

doc:
	go run cmd/kslibdocgen/main.go --path tmp/k8s.libsonnet
