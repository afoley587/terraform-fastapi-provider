SHELL:=/usr/bin/env bash

VERSION=0.0.1
HOSTNAME=terraform-example.com
PROVIDER=garbage
PLUGIN_DIR=${HOME}/.terraform.d/plugins
# For Darwin
CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wbuiltin-requires-header"

.PHONY: build lint

build:
	echo "building"
	CGO_CPPFLAGS=$(CGO_CPPFLAGS) go build -o terraform-provider-garbage
	mkdir -p ${PLUGIN_DIR}/${HOSTNAME}/${PROVIDER}/${PROVIDER}/${VERSION}/darwin_amd64/
	echo "copying"
	cp terraform-provider-garbage ${PLUGIN_DIR}/${HOSTNAME}/${PROVIDER}/${PROVIDER}/${VERSION}/darwin_amd64/

lint:
	go fmt
	go mod tidy