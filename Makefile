########################################################################################################################
# Copyright (c) 2019 meshbox Foundation
# This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
# warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
# permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
# License 2.0 that can be found in the LICENSE file.
########################################################################################################################

# Go parameters
GOCMD=go
GOLINT=golint
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOPATH=$(shell go env GOPATH)
BUILD_TARGET_SERVER=server

BUILD_TARGET_OS=$(shell go env GOOS)
BUILD_TARGET_ARCH=$(shell go env GOARCH)

# Pkgs
ALL_PKGS := $(shell go list ./... )
PKGS := $(shell go list ./... | grep -v /test/ )
ROOT_PKG := "github.com/MeshBoxFoundation/blockchain-node"

default: clean build
all: clean build

.PHONY: build
build:
	$(GOBUILD) -ldflags "$(PackageFlags)" -o ./bin/$(BUILD_TARGET_SERVER) -v ./src/blockchain_node_app.go

.PHONY: clean
clean:
	@echo "Cleaning..."
	$(ECHO_V)rm -rf ./bin/$(BUILD_TARGET_SERVER)

