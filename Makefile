# set USE_VARS=1 to use environment variables from the variables.env file
ifeq ($(USE_VARS), 1)
	include ./variables.env
endif

default: install

REPODIR=/tmp/tf-repo/providers

NAME=shoreline
BINARY=terraform-provider-$(NAME)
VERSION=1.15.39

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

version:
	@echo "version: ${VERSION}\nTo create a release run: \n  git tag v${VERSION}\n  git push origin v${VERSION}"
	@git tag -l | grep '^v${VERSION}$$' >/dev/null && echo "WARNING: Release already exists" || true

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


EXAMPLES_ROOT_PATH=./examples/resources/_root

init_ex:
	SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN=$(SHORELINE_TOKEN) SHORELINE_DEBUG=$(SHORELINE_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) init

apply_ex:
	SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN=$(SHORELINE_TOKEN) SHORELINE_DEBUG=$(SHORELINE_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) apply --auto-approve

apply_ex_na:
	SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN=$(SHORELINE_TOKEN) SHORELINE_DEBUG=$(SHORELINE_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) apply

destroy_ex:
	SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN=$(SHORELINE_TOKEN) SHORELINE_DEBUG=$(SHORELINE_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) destroy --auto-approve

destroy_ex_na:
	SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN=$(SHORELINE_TOKEN) SHORELINE_DEBUG=$(SHORELINE_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) destroy

plan_ex:
	SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN=$(SHORELINE_TOKEN) SHORELINE_DEBUG=$(SHORELINE_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) plan

.PHONY: distclean_ex
distclean_ex:
	rm -rf $(EXAMPLES_ROOT_PATH)/terraform.tfstate $(EXAMPLES_ROOT_PATH)/terraform.tfstate.backup $(EXAMPLES_ROOT_PATH)/.terraform $(EXAMPLES_ROOT_PATH)/.terraform.lock.hcl
