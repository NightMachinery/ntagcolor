# Build

`ntagcolor` is a normal Go command. Plain `go build` and `go install` work
because the generated runtime style table is checked in.

The recommended local build path is:

```sh
make build
```

`make build` runs `go generate ./...` first, then builds `./ntagcolor`. This
keeps `styles_gen.go` current when `styles_decl.go` changes.

## Targets

```sh
make generate  # refresh styles_gen.go
make build     # regenerate, then build ./ntagcolor
make test      # regenerate, then run go test ./...
make bench     # run Go benchmarks
make check     # regenerate, test, and fail if styles_gen.go changed
make clean     # remove ./ntagcolor
```

Use `make check` before committing changes to the color table or renderer. It
catches stale generated output while still allowing normal `go build` consumers
to install without running a custom build wrapper.

## Generated Styles

`styles_decl.go` is the editable color source. `go generate` runs
`generate_styles.go` and writes `styles_gen.go`, which contains resolved
foreground colors, background colors, and ANSI prefixes.

Commit both the declarative change and the regenerated `styles_gen.go`.
