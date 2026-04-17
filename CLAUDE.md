# code-cat

Go CLI tool (`ccat`) that automates Git workflows. Entry point: `cmd/ccat/main.go`.

## Testing

**After every bug fix or new feature**, spawn a sub-agent (Task tool) to write tests:

> "Write Go tests for the changes in [file(s)]. Use the standard `testing` package.
> Mirror patterns in the nearest existing `*_test.go` file. Run `go test -race -count=1 ./...`
> to verify all tests pass before finishing."

The sub-agent receives only the modified files — not the full conversation context.

Test command: `go test -race -count=1 ./...`
