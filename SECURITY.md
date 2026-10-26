# Security Policy

## Supported versions

Security fixes are provided for the latest released minor version. This
development branch is not a supported release until a version is published.

## Reporting a vulnerability

Please use
[GitHub private vulnerability reporting](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/security/advisories/new).
Do not open a public issue and do
not include credentials, access tokens, personal data, or production
configuration in a report.

Include the affected version, impact, reproduction steps, and a minimal
redacted example where possible. Maintainers will acknowledge the report,
coordinate remediation, and disclose it after a fix is available.

## Credential handling

Never commit or log Alibaba Cloud access keys, secret keys, STS credentials,
or IDaaS access tokens. Revoke and rotate any credential that may have been
exposed.
