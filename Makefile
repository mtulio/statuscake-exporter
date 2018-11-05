
export PATH := ./bin:$(PATH)
export GO111MODULE := on
export BIN_NAME := ./bin/statuscake-exporter

# Initial development
.PHONY: init
init:
	go mod init

# Install all the build and lint dependencies
.PHONY: setup
setup:
	go mod download

.PHONY: update
update:
	go mod tidy

# Build a beta version
.PHONY: build
build:
	@test -d ./bin || mkdir ./bin
	go build -o $(BIN_NAME) && strip $(BIN_NAME)

.PHONY: run
run:
	$(BIN_NAME)

.PHONY: version
version: build
	$(BIN_NAME) -v

.PHONY: clean
clean:
	@rm -f bin/$(BIN_NAME)

# ##########
# Goreleaser
# https://goreleaser.com/introduction/
gorelease-init:
	goreleaser init

# #######
# Release
tag:
	$(call deps_tag,$@)
	git tag -a $(shell cat VERSION) -m "$(message)"
	git push origin $(shell cat VERSION)

# Release tool
# https://goreleaser.com/introduction/
release:
	goreleaser --rm-dist