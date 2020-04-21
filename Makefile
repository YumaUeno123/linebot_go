export LINEBOT_HOME := $(PWD)
export GOBIN := $(PWD)/bin
export LINEBOT_VENDOR_PATH=${LINEBOT_HOME}/vendor

install-go:
	goenv install -s $$(cat .go-version)

install-modules:
	go mod tidy

install-tools:
	mkdir -p bin
	go install \
	golang.org/x/tools/cmd/goimports \
	github.com/golangci/golangci-lint/cmd/golangci-lint

install: install-go install-modules install-tools

goimports:
	$(GOBIN)/goimports -w $$(find . -name "*.go")

lint:
	$(GOBIN)/golangci-lint run