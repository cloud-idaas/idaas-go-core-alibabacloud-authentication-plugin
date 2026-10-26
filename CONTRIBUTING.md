# Contributing

Thank you for contributing.

## Development

This module supports Go 1.18 or later; newer toolchains are recommended
for development. Dependencies are resolved from their published Go modules.
CI verifies the 1.18 floor.

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
