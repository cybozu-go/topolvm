CSI_VERSION=1.1.0
PROTOBUF_VERSION=3.7.1
CURL=curl -Lsf

GOFLAGS = -mod=vendor
export GOFLAGS

PTYPES_PKG := github.com/golang/protobuf/ptypes
GO_OUT := plugins=grpc
GO_OUT := $(GO_OUT),Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor
GO_OUT := $(GO_OUT),Mgoogle/protobuf/wrappers.proto=$(PTYPES_PKG)/wrappers
GO_OUT := $(GO_OUT):csi

csi.proto:
	$(CURL) -o $@ https://raw.githubusercontent.com/container-storage-interface/spec/v$(CSI_VERSION)/csi.proto

bin/protoc:
	$(CURL) -o protobuf.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOBUF_VERSION)/protoc-$(PROTOBUF_VERSION)-linux-x86_64.zip
	unzip protobuf.zip bin/protoc 'include/*'
	rm -f protobuf.zip

bin/protoc-gen-go:
	GOBIN=$(shell pwd)/bin go install -mod=vendor github.com/golang/protobuf/protoc-gen-go

csi/csi.pb.go: csi.proto bin/protoc bin/protoc-gen-go
	mkdir -p csi
	PATH=$(shell pwd)/bin:$(PATH) bin/protoc -I. --go_out=$(GO_OUT) $<

lvmd/proto/lvmd.pb.go: lvmd/proto/lvmd.proto bin/protoc bin/protoc-gen-go
	PATH=$(shell pwd)/bin:$(PATH) bin/protoc -I. --go_out=plugins=grpc:. $<

test:
	test -z "$$(gofmt -s -l . | grep -v '^vendor' | tee /dev/stderr)"
	test -z "$$(golint $$(go list ./... | grep -v /vendor/) | tee /dev/stderr)"
	go install ./...
	go test -race -v ./...
	go vet ./...

.PHONY: test
