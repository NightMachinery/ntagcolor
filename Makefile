GO ?= go
BIN ?= ntagcolor

.PHONY: all generate build test bench check clean

all: build

generate:
	$(GO) generate ./...

build: generate
	$(GO) build -o $(BIN) .

test: generate
	$(GO) test ./...

bench:
	$(GO) test -bench=. -benchmem -run '^$$'

check: generate test
	git diff --exit-code -- styles_gen.go

clean:
	rm -f $(BIN)
