VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
PROVIDER_DIR := registry.opentofu.org/Berndinox/virtuacloud

.PHONY: build install test testacc docs generate lint

build:
	go build -o ./bin/terraform-provider-virtuacloud_$(VERSION) .

install: build
	mkdir -p ~/.terraform.d/plugins/$(PROVIDER_DIR)/$(VERSION)/$(shell go env GOOS)_$(shell go env GOARCH)
	cp ./bin/terraform-provider-virtuacloud_$(VERSION) ~/.terraform.d/plugins/$(PROVIDER_DIR)/$(VERSION)/$(shell go env GOOS)_$(shell go env GOARCH)/terraform-provider-virtuacloud_v$(VERSION)

test:
	go test ./internal/... -v -count=1

testacc:
	TF_ACC=1 go test ./internal/... -v -count=1 -timeout 120m

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

lint:
	golangci-lint run ./...

generate:
	go generate ./...

fmt:
	gofmt -w .