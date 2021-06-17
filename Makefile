default: install

PLUGDIR=~/.terraform.d/plugins
REPODIR=/tmp/tf-repo/providers
#SUBPATH=shoreline.io/terraform/shoreline/1.0.0/linux_amd64
SUBPATH=shoreline.io/terraform/shoreline/1.0.0/darwin_amd64

install:
	go build
	#mkdir -p $(PLUGDIR)/$(SUBPATH)
	#cp tf-json $(PLUGDIR)/$(SUBPATH)/
	rm -rf $(REPODIR)/*
	mkdir -p $(REPODIR)/$(SUBPATH)
	cp terraform-provider-shoreline $(REPODIR)/$(SUBPATH)/terraform-provider-shoreline

# Run acceptance tests
#############
# NOTE: SHORELINE_URL and SHORELINE_TOKEN should be set externally, e.g. 
#   SHORELINE_URL=https://test.us.api.shoreline-vm10.io
#   SHORELINE_TOKEN=xyz1lkajsdf8.kjalksdjfl...

.PHONY: testacc
testacc:
	@ SHORELINE_URL=$(SHORELINE_URL) SHORELINE_TOKEN="$(SHORELINE_TOKEN)" TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120s

# no checked in files should contain tokens
scan:
	find . -type f | xargs grep -l -e '[e]yJhb' || echo "scan is clean"
