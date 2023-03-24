default: install

REPODIR=/tmp/tf-repo/providers

NAME=shoreline
BINARY=terraform-provider-$(NAME)
VERSION=1.11.1

// NOTE: this only works for 64 bit linux and MacOs ("darwin")
OS=$(shell uname | tr 'A-Z' 'a-z')
#SUBPATH=shoreline.io/terraform/shoreline/$(VERSION)/$(OS)_amd64
SUBPATH=shorelinesoftware/local/shoreline/$(VERSION)/$(OS)_amd64

build: format
	go build
	go generate

test:
	echo unit-tests...

check:
	gofmt -l .

format:
	gofmt -w .

# NOTE: This relies on your ~/.terraformrc pointing to /tmp/tf-repo.
#   See terraformrc in the current dir
install: build
	rm -rf $(REPODIR)/*
	mkdir -p $(REPODIR)/$(SUBPATH)
	cp terraform-provider-shoreline $(REPODIR)/$(SUBPATH)/terraform-provider-shoreline

# This sets up your ~/.terraformrc (NOTE: need to re-run when the version changes)
use_local:
	@echo 'Setting up local overrides for shoreline provider in ~/.terraformrc'
	@echo 'NOTE: You need to re-run "make use_local" when the version changes."'
	@echo 'provider_installation { dev_overrides { "shorelinesoftware/shoreline" = "$(REPODIR)/$(SUBPATH)" } }' > ${HOME}/.terraformrc

use_registry:
	@echo 'Removing ~/.terraformrc, to use the terraform registry again (https://registry.terraform.io/providers/shorelinesoftware/shoreline/latest)'
	@rm ${HOME}/.terraformrc

release:
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_darwin_amd64
	GOOS=darwin  GOARCH=arm64 go build -o ./bin/$(BINARY)_$(VERSION)_darwin_arm64
	GOOS=linux   GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_linux_amd64
	GOOS=linux   GOARCH=arm64 go build -o ./bin/$(BINARY)_$(VERSION)_linux_arm64
	GOOS=linux   GOARCH=arm   go build -o ./bin/$(BINARY)_$(VERSION)_linux_arm
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_openbsd_amd64
	GOOS=windows GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_windows_amd64

# Run acceptance tests
#############
# NOTE: SHORELINE_URL and SHORELINE_TOKEN should be set externally, e.g.
#   SHORELINE_URL=https://test.us.api.shoreline-vm10.io
#   SHORELINE_TOKEN=xyz1lkajsdf8.kjalksdjfl...

.PHONY: testacc
testacc:
	@ SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN="$(SHORELINE_TOKEN)" TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 600s

# no checked in files should contain tokens
scan:
	find . -type f | xargs grep -l -e '[e]yJhb' || echo "scan is clean"
