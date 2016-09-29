# Set the build version
ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git rev-parse --short HEAD)
endif

# Setup some useful vars
HOST_GOOS = $(shell go env GOOS)
HOST_GOARCH = $(shell go env GOARCH)
GLIDE_VERSION = v0.11.1
ifeq ($(origin GLIDE_GOOS), undefined)
	GLIDE_GOOS := $(HOST_GOOS)
endif

build: vendor
	go build -o bin/kuberang -ldflags "-X main.version=$(VERSION)" ./cmd

clean:
	rm -rf bin
	rm -rf out
	rm -rf vendor

test: vendor
	go test ./cmd/... ./pkg/... $(TEST_OPTS)

integration-test: dist just-integration-test

vendor: tools/glide
	./tools/glide install

tools/glide:
	mkdir -p tools
	curl -L https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-$(GLIDE_GOOS)-$(HOST_GOARCH).tar.gz | tar -xz -C tools
	mv tools/$(GLIDE_GOOS)-$(HOST_GOARCH)/glide tools/glide
	rm -r tools/$(GLIDE_GOOS)-$(HOST_GOARCH)

dist: build
	mkdir -p out
	cp bin/kuberang out