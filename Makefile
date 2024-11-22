LOCAL_BIN:=$(CURDIR)/bin
IMAGE_PREVIEWER_APP:=$(LOCAL_BIN)/image-previewer
DOCKER_IMG="image-previewer:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(IMAGE_PREVIEWER_APP) -ldflags "$(LDFLAGS)" -mod vendor ./cmd/image-previewer/.

run: build
	$(IMAGE_PREVIEWER_APP)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f $(CURDIR)/build/Dockerfile .

run-img: build-img
	docker run -p 8080:8080 $(DOCKER_IMG)

test:
	go test --race -count 100 ./... -v

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.0

lint: install-lint-deps
	golangci-lint run

.PHONY: build run build-img run-img test install-lint-deps lint