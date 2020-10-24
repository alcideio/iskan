#
# iskan --> Locates : Locate k8s, Helm & kustomize configuration issues and provide recommendation
#
.SECONDARY:
.SECONDEXPANSION:

BINDIR      := $(CURDIR)/bin
DIST_DIRS   := find * -type d -exec
DIST_EXES   := find * -type f -executable -exec
BINNAME     ?= iskan

GOPATH        = $(shell go env GOPATH)
DEP           = $(GOPATH)/bin/dep
GOX           = $(GOPATH)/bin/gox
GOIMPORTS     = $(GOPATH)/bin/goimports
ARCH          = $(shell uname -p)

UPX_VERSION := 3.96
UPX := $(CURDIR)/iskan/bin/upx

GORELEASER_VERSION := 0.145.0
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

#
# BUILD TOOL CHAIN
#

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

#
# BUILD
#

.PHONY: build
build: ##@build Build on local platform
	export CGO_ENABLED=0 && go build -o $(BINDIR)/$(BINNAME) -tags staticbinary -v -ldflags '$(LDFLAGS)'  github.com/alcideio/iskan

#
# Test & E2E
#

.PHONY: test
test: ##@test run tests
	go test -v github.com/alcideio/iskan/pkg/...

.PHONY: e2e
e2e-build: ##@test run tests
	# For local test runs you need to set E2E_GCR_PULLSECRET and E2E_API_CONFIG - see e2e/framework/config.go
	go test -v github.com/alcideio/iskan/e2e  -c -o bin/e2e.test

#
# To run specific e2e spec
# bin/e2e.test -v 8 -ginkgo.v -ginkgo.focus="\[local-scanner\]"
#
e2e: e2e-build ##@test run tests
	# For local test runs you need to set E2E_GCR_PULLSECRET and E2E_API_CONFIG - see e2e/framework/config.go
	e2e/e2e-pipeline-runner.sh


e2e-coverage: ##@test run tests
	go test -v -race -covermode atomic -coverprofile=e2e-coverage.out github.com/alcideio/iskan/e2e
	go tool cover -func=e2e-coverage.out
	go tool cover -html=e2e-coverage.out -o e2e-coverage.html

.PHONY: coverage
coverage: ##@test run tests with coverage report
	go test -v -race -covermode atomic -coverprofile=coverage.out github.com/alcideio/iskan/pkg/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

#
# Make sure to login to those registries
# 1. GCR: gcloud auth configure-docker && gcloud auth login
# 2. ECR: aws --profile iskan ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 893825821121.dkr.ecr.us-west-2.amazonaws.com
# 3. ACR: az acr login --name alcide --subscription 9efc9618-47a0-4e98-b31e-7194f25188d4 --resource-group  iskan
# 4. Docker Hub: docker login --username alcide --password <GET PASSWORD REPO>
#
REGISTRIES ?= 893825821121.dkr.ecr.us-west-2.amazonaws.com/iskan gcr.io/dcvisor-162009/iskan/e2e alcide.azurecr.io/iskan iskan
e2e-build-images: ##@e2e build test images
	@for reg in ${REGISTRIES}; \
	do \
		for dir in $(wildcard e2e/images/*); \
		do \
			  image=`basename $$dir` ; cd $$dir ; echo $$reg/$$image > build.txt ;docker build -t $$reg/$$image . ; docker push $$reg/$$image ; cd -;\
		done \
    done

create-kind-cluster:  ##@Test creatte KIND cluster
	kind create cluster --image kindest/node:v1.18.2 --name iskan

delete-kind-cluster:  ##@Test delete KIND cluster
	kind delete cluster --name iskan

#
# Helm
#
helm-docs:  ##@Helm Generate Documentation
	docker run --rm --volume $(CURDIR):/helm-docs jnorwood/helm-docs:latest

helm-lint:  ##@Helm Generate Documentation
	helm lint deploy/charts/iskan

helm-gen:  ##@Helm Generate Kubernetes manifest
	helm template -n alcide-iskan iskan deploy/charts/iskan

helm-install:  ##@Helm Install
	helm upgrade -i -n alcide-iskan iskan deploy/charts/iskan

helm-delete:  ##@Helm Delete Installation
	helm -n alcide-iskan delete iskan

#
# HTML Viewer
#
viewer-deps:  ##@HTMLViewer Launch Dev Server for development
	cd htmlviewer && npm install

viewer-dev:  ##@HTMLViewer Launch Dev Server for development
	cd htmlviewer && npm run dev

#
# 1. Get alcide builder NPMJS token from 1password
# 2. export NODE_AUTH_TOKEN=<The Token>
viewer-release:  ##@HTMLViewer Launch Dev Server for development
	cd htmlviewer &&\
	npm clean-install &&\
	npm run prod &&\
	npm publish --access public
#
# RELEASE
#

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
helm-release:  ##@release create helm chart archive
	mkdir artifacts || true
	cd deploy && tar --strip-components=2 -cvzf ../artifacts/iskan-helm-chart.tar.gz charts

.PHONY: gorelease
gorelease: helm-release ##@release Generate All release artifacts
	GOPATH=~ USER=alcidebuilder $(GORELEASER) -f $(CURDIR)/.goreleaser.yml --rm-dist --release-notes=notes.md

gorelease-snapshot: helm-release ##@release Generate All release artifacts
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
