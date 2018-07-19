.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	go get -u github.com/golang/dep/cmd/dep

.PHONY: dep
dep: ## Run dep ensure
	dep ensure

.PHONY: lint
lint: ## Run all the linters
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=goconst \
		--enable=goimports \
		--enable=gosimple \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=maligned \
		--enable=misspell \
		--enable=staticcheck \
		--enable=unconvert \
		--enable=varcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...

.PHONY: build
build: ## Build a binary
	go build -v -o ./bin/outlyer ./cmd/outlyer

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
