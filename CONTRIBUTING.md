# Contributing

Thank you for contributing.

## Development

Use Go 1.21 or later. Dependencies are resolved from their published Go
modules.

Before submitting a change, run:

```bash
gofmt -w .
go vet ./...
go test -race -cover ./...
```

Add tests for behavior changes and avoid committing credentials, tokens,
generated coverage files, or environment-specific configuration. Keep public
API changes backward compatible unless a breaking release has been agreed.

By participating, you agree to follow [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
