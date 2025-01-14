BUILD_BASE=$(shell pwd)
GOBASE=$(BUILD_BASE)
GOBIN=$(GOBASE)/bin
binary_name=go-nordvpn
binary=$(GOBIN)/$(binary_name)
source=cmd/$(binary_name).go
image_name=aharper343/$(binary_name)


help:
	@echo "This is a helper makefile for go-nordvpn"
	@echo "Targets:"
	@echo "    generate:    regenerate all generated files"
	@echo "    test:        run all tests"
	@echo "    tidy:        tidy go mod"
	@echo "    lint:        lint the project"
	@echo "    all:         runs clean generate build on the project"
	@echo "    build:       builds the project binary"
	@echo "    run:         runs the project binary"
	@echo "    pre-commit:  runs tidy, verify, format and lint on the project"
	@echo "    verify:      verifies the project"
	@echo "    format:      formats the source files"

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.62.2

.PHONY: tools
tools: $(GOBIN)/golangci-lint

generate:
	go generate ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

tidy-ci:
	tidied -verbose

lint: tools
	$(GOBIN)/golangci-lint run ./...

lint-ci: tools
	$(GOBIN)/golangci-lint run ./... --out-format=colored-line-number --timeout=5m

all: clean generate build

clean:
	go clean ./...
	rm -f */*.gen.go $(binary)

build:
	go build -o $(binary) $(source)

run:
	go run $(source)

docker-build:
	docker build --tag $(image_name):$$(date +v%Y.%m.%d-%H%M%S) --tag $(image_name):latest .

pre-commit: tidy verify format lint

verify:
	go mod verify

format:
	go fmt ./...
