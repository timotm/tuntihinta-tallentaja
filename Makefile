
GO?=go

all: tuntihinta-tallentaja test

.PHONY: test
test:
	$(GO) test ./...

tuntihinta-tallentaja: $(shell git ls-files "*.go")
	$(GO) build ./cmd/tuntihinta-tallentaja
