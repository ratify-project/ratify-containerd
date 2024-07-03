WATCHER_BINARY_NAME = configmap-watcher

all: build

.PHONY: build
build: build-watcher

.PHONY: build-watcher
build-watcher: fmt vet
	go build -o ./bin/${WATCHER_BINARY_NAME}

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...