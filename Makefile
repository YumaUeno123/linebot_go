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