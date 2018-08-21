#!/usr/bin/make -f

export CGO_ENABLED=0

VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-list -1 HEAD)

define gox_build
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/$(1)_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static" -X github.com/previousnext/mysql-toolkit/internal/version.GitVersion=${VERSION} -X github.com/previousnext/mysql-toolkit/internal/version.GitCommit=${COMMIT}' \
        $(2)
endef

# Build all.
build: cli operator

# Builds the main project.
cli:
	$(call gox_build,mtk,github.com/previousnext/mysql-toolkit/cmd/mtk)

# Build the Kubernetes operator.
operator:
	$(call gox_build,mtk-operator,github.com/previousnext/mysql-toolkit/cmd/mtk-operator)

# Run all lint checking with exit codes for CI
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test:
	go test -cover ./...

# Release the project to all.
release: release-docker

define docker_build_push
	docker build -f $(1) -t $(2) .
	docker push $(2)
endef

# Release the project to Docker Registry.
release-docker:
	$(call docker_build_push,containers/cli/Dockerfile,previousnext/mtk:$(VERSION))
	$(call docker_build_push,containers/operator/Dockerfile,previousnext/mtk-operator:$(VERSION))

.PHONY: build cli operator lint test release release-docker
