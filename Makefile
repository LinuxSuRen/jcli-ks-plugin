NAME := jcli-ks-plugin
CGO_ENABLED = 0
BUILD_GOOS=$(shell go env GOOS)
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILDFLAGS = -ldflags "-X github.com/linuxsuren/cobra-extension/version.version=$(VERSION) \
	-X github.com/linuxsuren/cobra-extension/version.commit=$(COMMIT) \
	-X github.com/linuxsuren/cobra-extension/version.date=$(shell date +'%Y-%m-%d')"
COVERED_MAIN_SRC_FILE=./main
PATH := $(PATH):$(PWD)/bin

build: fmt
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=$(BUILD_GOOS) GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) \
	-o bin/$(BUILD_GOOS)/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/$(BUILD_GOOS)/$(NAME)

run:
	go run main.go

fmt:
	go fmt ./...

copy: build
	cp bin/$(BUILD_GOOS)/$(NAME) ~/.jcli-plugins/plugins-repo/$(NAME)

