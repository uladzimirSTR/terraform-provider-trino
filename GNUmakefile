PROVIDER_NAME := trino
PROVIDER_NAMESPACE := $(shell whoami)
PROVIDER_HOSTNAME := registry.terraform.io
VERSION := 0.1.0
OS_ARCH := linux_amd64

BINARY := terraform-provider-$(PROVIDER_NAME)
LOCAL_PLUGIN_DIR := ~/.terraform.d/plugins/$(PROVIDER_HOSTNAME)/$(PROVIDER_NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS_ARCH)

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o $(BINARY)

.PHONY: install-local
install-local: build
	mkdir -p $(LOCAL_PLUGIN_DIR)
	mv $(BINARY) $(LOCAL_PLUGIN_DIR)/$(BINARY)

.PHONY: clean
clean:
	rm -f $(BINARY)

tfplugindocs:
	export GOBIN=$PWD/bin
	export PATH=$GOBIN:$PATH
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

docs: tfplugindocs
	tfplugindocs generate

.PHONY: example-init
example-init: install-local
	cd examples/schema && terraform init

.PHONY: example-plan
example-plan: install-local
	cd examples/schema && terraform plan

.PHONY: example-apply
example-apply: install-local
	cd examples/schema && terraform apply

.PHONY: example-destroy
example-destroy: install-local
	cd examples/schema && terraform destroy
