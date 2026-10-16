# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

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
