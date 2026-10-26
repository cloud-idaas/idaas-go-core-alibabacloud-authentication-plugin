# idaas-go-core-alibabacloud-authentication-plugin

[![CI](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/actions/workflows/ci.yml/badge.svg)](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.18%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Development Status](https://img.shields.io/badge/status-Beta-yellow)](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin)

[简体中文](README_zh.md) | English


Authentication plugin for machine-to-machine access to Alibaba Cloud IDaaS
EIAM. The plugin obtains an Alibaba Cloud credential from the official default
credential chain and calls EIAM
[GenerateOauthToken](https://www.alibabacloud.com/help/en/idaas/eiam/developer-reference/api-eiam-2021-12-01-generateoauthtoken)
to acquire an OAuth2 access token.


## Features

- **Alibaba Cloud EIAM Integration**: Acquires OAuth2 access tokens via the EIAM
  `GenerateOauthToken` API using the Alibaba Cloud default credential chain.
- **Plugin Architecture**: Implements the `PluginCredentialProvider` interface
  from `idaas-go-core-sdk`.
- **Automatic Registration**: Auto-registers with the Core SDK plugin registry
  via blank import — no manual setup required.

## Overview

```
┌─────────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│   idaas-go-core-sdk  │────▶│  this plugin     │────▶│  EIAM OpenAPI    │
│  (config · factory  │     │  (Alibaba Cloud   │     │  GenerateOauthToken│
│   · cache · refresh)│     │   credential →    │     │  → OAuth2 token  │
└─────────────────────┘     │   EIAM token)     │     └──────────────────┘
                            └──────────────────┘
```

The Core SDK loads configuration and coordinates token caching/refresh;
this plugin obtains the Alibaba Cloud credential and calls EIAM
`GenerateOauthToken`.

## Requirements

- Go 1.18 or later, matching the minimum version declared by
  [`idaas-go-core-sdk`](https://github.com/cloud-idaas/idaas-go-core-sdk).
  Note: Go 1.18 is end-of-life upstream and only used as the compatibility
  floor verified by CI; prefer a supported Go release for new deployments
- An Alibaba Cloud credential accepted by
  [credentials-go](https://github.com/aliyun/credentials-go)
  (environment variables, instance RAM role, etc.)
- An IDaaS EIAM application authorized for the requested scopes

### Dependencies

| Module | Minimum version |
|---|---|
| [`idaas-go-core-sdk`](https://github.com/cloud-idaas/idaas-go-core-sdk) | v0.1.0-beta.1 |
| [`alibabacloud-go/eiam-20211201/v2`](https://github.com/alibabacloud-go/eiam-20211201) | v2.15.9 |
| [`credentials-go`](https://github.com/aliyun/credentials-go) | v1.4.5 |

## Installation

```bash
go get github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin
```

## Configuration

The plugin name is stable and must be configured as
`alibabacloudPluginCredentialProvider`.

Create `~/.cloud_idaas/client-config.json`:

```json
{
  "idaasInstanceId": "idaas_xxx",
  "clientId": "app_xxx",
  "issuer": "https://example.com/api/v2/iauths_system/oauth2",
  "tokenEndpoint": "https://example.com/api/v2/iauths_system/oauth2/token",
  "scope": "api.example.com|read:file api.example.com|write:file",
  "openApiEndpoint": "eiam.cn-hangzhou.aliyuncs.com",
  "authnConfiguration": {
    "identityType": "CLIENT",
    "authnMethod": "PLUGIN",
    "pluginName": "alibabacloudPluginCredentialProvider"
  },
  "httpConfiguration": {
    "connectTimeout": 5000,
    "readTimeout": 10000
  }
}
```

### Plugin-specific fields

| Field | Required | Description |
|---|---:|---|
| idaasInstanceId | Yes | IDaaS EIAM instance ID |
| clientId | Yes | IDaaS application ID |
| scope | Yes | One or more `audience\|scopeValue` pairs (max 4096 chars, 100 values, single audience) |
| openApiEndpoint | Yes | Official EIAM endpoint (`eiam.<region>.aliyuncs.com` or `eiam-vpc.<region>.aliyuncs.com`). A bare host is accepted and normalized to `https://` at runtime; an explicit scheme must be `https` |
| identityType | Yes | Must be `CLIENT` |
| authnMethod | Yes | Must be `PLUGIN` (the Go Core SDK enum value, case-insensitive) |
| pluginName | Yes | Must be `alibabacloudPluginCredentialProvider` |

### Core SDK shared fields

| Field | Description |
|---|---|
| issuer | OIDC issuer endpoint for the IDaaS instance |
| tokenEndpoint | OAuth2 token endpoint URL |

| Field | Default | Valid range | Description |
|---|---:|---|---|
| httpConfiguration.connectTimeout | 5000 | 2000–60000 ms | TCP connection timeout |
| httpConfiguration.readTimeout | 10000 | 2000–60000 ms | HTTP read timeout |

Timeouts are in milliseconds. `0` inherits the default; out-of-range values are
rejected at startup.

## Usage

Import the plugin for its registration side effect:

```go
package main

import (
    "fmt"
    "log"

    "github.com/cloud-idaas/idaas-go-core-sdk/config"
    "github.com/cloud-idaas/idaas-go-core-sdk/factory"
    _ "github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin"
)

func main() {
    cfg, err := config.NewConfigReader().LoadWithPriority(
        "~/.cloud_idaas/client-config.json",
    )
    if err != nil {
        log.Fatal(err)
    }

    coreFactory := factory.GetInstance()
    if err := coreFactory.Initialize(cfg); err != nil {
        log.Fatal(err)
    }

    provider, err := coreFactory.CreateCredentialProvider()
    if err != nil {
        log.Fatal(err)
    }
    cred, err := provider.GetCredential()
    if err != nil {
        log.Fatal(err)
    }

    // Never print or log the access token.
    fmt.Printf("Credential obtained: type=%s expiresIn=%ds\n",
        cred.GetTokenType(), cred.GetExpiresIn())
}
```

Run the included sample:

```bash
go run ./samples/plugin_authentication.go \
  -config ~/.cloud_idaas/client-config.json
```

## Troubleshooting

### Configuration errors

| Error code | Cause |
|---|---|
| `IDaaSCredentialProviderFactoryNotInit` | Factory not initialized — call `factory.GetInstance().Initialize(cfg)` first |
| `IDaaSInstanceIdNotFound` | `idaasInstanceId` is empty |
| `ClientIdNotFound` | `clientId` is empty |
| `OpenApiEndpointNotFound` | `openApiEndpoint` is empty, not HTTPS, not an official EIAM endpoint, or contains a path/IP/non-443 port |
| `ConnectTimeoutNotValid` | `connectTimeout` is outside 2000–60000 ms |
| `ReadTimeoutNotValid` | `readTimeout` is outside 2000–60000 ms |
| `InvalidScope` | Scope is empty, exceeds 4096 chars, has more than 100 values, or is malformed (not `audience\|scopeValue`) |
| `MultipleAudienceNotSupported` | Scope items reference more than one audience |

### Runtime errors

| Error code | Cause |
|---|---|
| `AlibabaCloudCredentialError` | The Alibaba Cloud default credential chain could not find a credential |
| `EIAMClientError` | EIAM SDK client creation failed |
| `GenerateOAuthTokenError` | The EIAM `GenerateOauthToken` API call failed |
| `InvalidTokenResponse` | EIAM returned an incomplete, empty, or expired token |

When the EIAM API returns an HTTP error, it is classified as follows:

| HTTP status | Error category |
|---|---|
| 500+ | ServerError |
| 400–499 | ClientError |
| Other (`0` = transport failure without an HTTP response, e.g. connect timeout or DNS error) | HttpError |

## Development

```bash
gofmt -w .
go vet ./...
go test -race -cover ./...
```

Integration tests against the real EIAM API: see
[INTEGRATION_TEST.md](INTEGRATION_TEST.md).

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), and
[CHANGELOG.md](CHANGELOG.md).

## Support and feedback

- **Issues**: Please open a
  [GitHub issue](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/issues)
  for questions or suggestions.
- **Security vulnerabilities**: see [SECURITY.md](SECURITY.md) — please do not
  open a public issue.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
