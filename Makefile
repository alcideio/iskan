#
# iskan --> Locates : Locate k8s, Helm & kustomize configuration issues and provide recommendation
#
.SECONDARY:
.SECONDEXPANSION:

BINDIR      := $(CURDIR)/bin
DIST_DIRS   := find * -type d -exec
DIST_EXES   := find * -type f -executable -exec
# Go Targets darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
TARGETS     := darwin/amd64 linux/amd64 windows/amd64
TARGET_OBJS ?= darwin-amd64.tar.gz darwin-amd64.tar.gz.sha256 darwin-amd64.tar.gz.sha256sum linux-amd64.tar.gz linux-amd64.tar.gz.sha256 linux-amd64.tar.gz.sha256sum linux-386.tar.gz linux-386.tar.gz.sha256 linux-386.tar.gz.sha256sum linux-arm.tar.gz linux-arm.tar.gz.sha256 linux-arm.tar.gz.sha256sum linux-arm64.tar.gz linux-arm64.tar.gz.sha256 linux-arm64.tar.gz.sha256sum linux-ppc64le.tar.gz linux-ppc64le.tar.gz.sha256 linux-ppc64le.tar.gz.sha256sum linux-s390x.tar.gz linux-s390x.tar.gz.sha256 linux-s390x.tar.gz.sha256sum windows-amd64.zip windows-amd64.zip.sha256 windows-amd64.zip.sha256sum
BINNAME     ?= iskan

GOPATH        = $(shell go env GOPATH)
DEP           = $(GOPATH)/bin/dep
GOX           = $(GOPATH)/bin/gox
GOIMPORTS     = $(GOPATH)/bin/goimports
ARCH          = $(shell uname -p)

UPX_VERSION := 3.96
UPX := $(CURDIR)/iskan/bin/upx

GORELEASER_VERSION := 0.131.1
GORELEASER := $(CURDIR)/bin/goreleaser

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

# Required for globs to work correctly
#SHELL      = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

LDFLAGS += -X github.com/alcideio/iskan/pkg/version.Commit=${GIT_SHA}

BINARY_VERSION ?= ${GIT_TAG}
ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif


## Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/alcideio/iskan/pkg/version.Version=${BINARY_VERSION}
endif

get-bins: get-release-bins ##@build Download UPX
	wget https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz && \
	tar xvf upx-${UPX_VERSION}-amd64_linux.tar.xz &&\
	mkdir -p $(CURDIR)/bin || echo "dir already exist" &&\
	cp upx-${UPX_VERSION}-amd64_linux/upx $(CURDIR)/bin/upx &&\
	rm -Rf upx-${UPX_VERSION}-amd64_linux*

get-release-bins: ##@build Download goreleaser
	mkdir -p $(CURDIR)/bin || echo "dir already exist" &&\
	cd $(CURDIR)/bin &&\
	wget https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_Linux_x86_64.tar.gz && \
	tar xvf goreleaser_Linux_x86_64.tar.gz &&\
	rm -Rf goreleaser_Linux_x86_64*


.PHONY: build
build: ##@build Build on local platform
	export CGO_ENABLED=0 && go build -o $(BINDIR)/$(BINNAME) -tags staticbinary -v -ldflags '$(LDFLAGS)'  github.com/alcideio/iskan

.PHONY: test
test: ##@test run tests
	go test -v github.com/alcideio/iskan/...

.PHONY: coverage
coverage: ##@test run tests with coverage report
	go test -v -race -covermode atomic -coverprofile=coverage.out github.com/alcideio/iskan/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

create-kind-cluster:  ##@Test creatte KIND cluster
	kind create cluster --image kindest/node:v1.18.2 --name iskan

delete-kind-cluster:  ##@Test delete KIND cluster
	kind delete cluster --name iskan

#
#  How to release:
#
#  1. Grab GITHUB Token of alcidebuilder from 1password
#  2. export GITHUB_TOKEN=<thetoken>
#  3. git tag -a v0.4.0 -m "my new version"
#  4. git push origin v0.4.0
#  5. Go to to https://github.com/alcideio/iskan/releases and publish the release draft
#
#  Delete tag: git push origin --delete v0.7.0
#
.PHONY: gorelease
gorelease: ##@build Generate All release artifacts
	GOPATH=~ USER=alcidebuilder $(GORELEASER) -f $(CURDIR)/.goreleaser.yml --rm-dist --release-notes=notes.md

gorelease-snapshot: ##@build Generate All release artifacts
	GOPATH=~ USER=alcidebuilder  GORELEASER_CURRENT_TAG=v0.0.0 $(GORELEASER) -f $(CURDIR)/.goreleaser.yml --rm-dist --skip-publish --snapshot

HELP_FUN = \
         %help; \
         while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^(.+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/ }; \
         print "Usage: make [opti@buildons] [target] ...\n\n"; \
     for (sort keys %help) { \
         print "$$_:\n"; \
         for (sort { $$a->[0] cmp $$b->[0] } @{$$help{$$_}}) { \
             $$sep = " " x (30 - length $$_->[0]); \
             print "  $$_->[0]$$sep$$_->[1]\n" ; \
         } print "\n"; }

help: ##@Misc Show this help
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
