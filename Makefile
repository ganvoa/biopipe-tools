.PHONY: build
build:
	@go build -o ./biopipe ./cmd/cli/biopipe.go

.PHONY: push
push:
	@docker push gamboa/biopipe-cli:latest

.PHONY: mocks
mocks:
	@rm -rf ./mocks
	@bin/mockery --all