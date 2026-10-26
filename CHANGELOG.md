# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [v0.1.0-beta.1] - 2026-10-26

Incremental changes on top of 0.1.0-alpha.1.

### Added

- Apache-2.0 license headers on every source file.
- Dependabot updates for Go modules and GitHub Actions.
- Runnable godoc example for the plugin registry.
- English integration-test guide (`INTEGRATION_TEST.md`) alongside the
  Chinese version (`INTEGRATION_TEST_zh.md`), both linked from the READMEs.
- CI step that builds all packages, including `samples/`.

### Changed

- Set the minimum supported Go version to 1.18 to match the Core SDK
  (CI verifies the 1.18 floor).
- Upgrade `idaas-go-core-sdk` dependency to v0.1.0-beta.1.
- Keep `openApiEndpoint` examples as bare hosts; the Core SDK normalizes
  them to `https://` before validation during `factory.Initialize`.
- Document `authnMethod` case handling: configuration examples keep the
  uppercase `PLUGIN` spelling returned by the IDaaS configuration API, while
  README usage examples show the canonical enum value `plugin`; the Core SDK
  normalizes the field case-insensitively on config load.
- Document that error responses carrying HTTP status code `0` indicate a
  transport failure without an HTTP response.
- Replace the internal email-list contact with GitHub Issues and
  GitHub private vulnerability reporting across all public documentation.

### Fixed

- Reject EIAM endpoint region labels with a trailing hyphen
  (e.g. `eiam.cn-hangzhou-.aliyuncs.com`).

## [0.1.0-alpha.1] - 2026-10-16

### Fixed

- Classify transport failures without an HTTP response as HTTP errors instead
  of 4xx client errors.
- Validate optional connect and read timeouts against the documented
  2,000-60,000 millisecond range before creating the EIAM client.

### Added

- Explicit provider construction API with automatic plugin registration.
- Strict scope, endpoint, and token response validation.
- Configurable client dependencies and clock for deterministic testing.
- Unit, fuzz, race, and error-mapping tests.

### Changed

- Align the module path with the `cloud-idaas` GitHub organization.
- Split provider responsibilities into focused source files.
- Propagate configured connect and read timeouts to EIAM requests.
- Return stable, classified errors for credential and service failures.

### Security

- Restrict endpoints to HTTPS Alibaba Cloud EIAM hosts and port 443.
- Remove token material from sample output.
- Add an actionable private vulnerability-reporting channel and CI scan.

[v0.1.0-beta.1]: https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/releases/tag/v0.1.0-beta.1
