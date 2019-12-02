
# ###########
# Global Vars

export PATH := ./bin:$(PATH)
export GO111MODULE := on
BIN_NAME := ./bin/statuscake-exporter

APP_NAME ?= statuscake-exporter
VERSION := $(shell cat ./VERSION)
ENV := production

DOCKER_REPO 		?= mtulio
DOCKER_IMAGE_NAME 	= $(APP_NAME)
DOCKER_IMAGE_TAG 	?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

CPWD := $(PWD)

TMP_DIRS := ./bin
TMP_DIRS += ./dist

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_DESCRIBE := $(shell git describe --tags --always)
GIT_BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)

BUILD_DATE := $(shell date +%Y/%m/%d-%H:%M:%S)

GOOS := linux
GOARCH := amd64

CGO_ENABLED := 0

LDFLAGS :=
LDFLAGS += -X github.com/prometheus/common/version.Branch=$(GIT_COMMIT)
LDFLAGS += -X github.com/prometheus/common/version.Revision==$(GIT_DESCRIBE)
LDFLAGS += -X github.com/prometheus/common/version.Version=$(VERSION)
LDFLAGS += -X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE)
LDFLAGS += -X main.VersionEnv=$(ENV)

GORELEASE_VERSION 	:= v0.86.1
GORELEASE_BASE_URL 	:= https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASE_VERSION)/goreleaser
GORELEASE_URL_RPM 	:= $(GORELEASE_BASE_URL)_amd64.rpm
