# idaas-go-core-alibabacloud-authentication-plugin

[![CI](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/actions/workflows/ci.yml/badge.svg)](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Development Status](https://img.shields.io/badge/status-Alpha-orange)](https://github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin)

[English](README.md) | 简体中文


本插件面向机器到机器认证，使用阿里云官方默认凭证链获取云凭证，并调用
IDaaS EIAM
[GenerateOauthToken](https://www.alibabacloud.com/help/en/idaas/eiam/developer-reference/api-eiam-2021-12-01-generateoauthtoken)
接口获取 OAuth2 访问令牌。


## 特性

- **阿里云 EIAM 集成**：通过 EIAM `GenerateOauthToken` 接口，使用阿里云默认凭证链获取 OAuth2 访问令牌。
- **插件架构**：实现 `idaas-go-core-sdk` 的 `PluginCredentialProvider` 接口。
- **自动注册**：通过匿名导入自动注册到 Core SDK 插件注册表 — 无需手动配置。

## 概述

```
┌─────────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│   idaas-go-core-sdk  │────▶│  本插件           │────▶│  EIAM OpenAPI    │
│  (配置 · 工厂        │     │  (阿里云凭证 →     │     │  GenerateOauthToken│
│   · 缓存 · 刷新)     │     │   EIAM 令牌)      │     │  → OAuth2 令牌   │
└─────────────────────┘     └──────────────────┘     └──────────────────┘
```

Core SDK 加载配置并协调令牌缓存/刷新；本插件获取阿里云凭证并调用
EIAM `GenerateOauthToken`。

## 环境要求

- Go 1.21 及以上
- credentials-go 支持的阿里云凭证
  （环境变量、实例 RAM 角色等）
- 已授权目标 Scope 的 IDaaS EIAM 应用

### 依赖项

| 模块 | 最低版本 |
|---|---|
| [`idaas-go-core-sdk`](https://github.com/cloud-idaas/idaas-go-core-sdk) | v0.1.0-alpha.3 |
| [`alibabacloud-go/eiam-20211201/v2`](https://github.com/alibabacloud-go/eiam-20211201) | v2.15.9 |
| [`credentials-go`](https://github.com/aliyun/credentials-go) | v1.4.5 |

## 安装

```bash
go get github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin
```

## 配置

插件名固定，必须配置为 `alibabacloudPluginCredentialProvider`。

创建 `~/.cloud_idaas/client-config.json`：

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
    "authnMethod": "plugin",
    "pluginName": "alibabacloudPluginCredentialProvider"
  },
  "httpConfiguration": {
    "connectTimeout": 5000,
    "readTimeout": 10000
  }
}
```

### 插件特有字段

| 字段 | 必填 | 说明 |
|---|---:|---|
| idaasInstanceId | 是 | IDaaS EIAM 实例 ID |
| clientId | 是 | IDaaS 应用 ID |
| scope | 是 | 一个或多个 `audience\|scopeValue` 对（最长 4096 字符，最多 100 个值，单一 audience） |
| openApiEndpoint | 是 | 官方 EIAM Endpoint（`eiam.<region>.aliyuncs.com` 或 `eiam-vpc.<region>.aliyuncs.com`） |
| identityType | 是 | 必须为 `CLIENT` |
| authnMethod | 是 | 必须为 `plugin`（Go Core SDK 枚举值） |
| pluginName | 是 | 必须为 `alibabacloudPluginCredentialProvider` |

### Core SDK 共享字段

| 字段 | 说明 |
|---|---|
| issuer | IDaaS 实例的 OIDC 签发者 Endpoint |
| tokenEndpoint | OAuth2 令牌 Endpoint URL |

| 字段 | 默认值 | 有效范围 | 说明 |
|---|---:|---|---|
| httpConfiguration.connectTimeout | 5000 | 2000–60000 ms | TCP 连接超时 |
| httpConfiguration.readTimeout | 10000 | 2000–60000 ms | HTTP 读取超时 |

超时单位为毫秒。`0` 使用默认值；超出范围的值会在启动时被拒绝。

## 使用

通过匿名导入完成自动注册：

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

    // 切勿打印或记录访问令牌。
    fmt.Printf("Credential obtained: type=%s expiresIn=%ds\n",
        cred.GetTokenType(), cred.GetExpiresIn())
}
```

运行示例：

```bash
go run ./samples/plugin_authentication.go \
  -config ~/.cloud_idaas/client-config.json
```

## 故障排查

### 配置错误

| 错误码 | 原因 |
|---|---|
| `IDaaSCredentialProviderFactoryNotInit` | 工厂未初始化 — 请先调用 `factory.GetInstance().Initialize(cfg)` |
| `IDaaSInstanceIdNotFound` | `idaasInstanceId` 为空 |
| `ClientIdNotFound` | `clientId` 为空 |
| `OpenApiEndpointNotFound` | `openApiEndpoint` 为空、非 HTTPS、非官方 EIAM Endpoint，或包含路径/IP/非 443 端口 |
| `ConnectTimeoutNotValid` | `connectTimeout` 不在 2000–60000 ms 范围内 |
| `ReadTimeoutNotValid` | `readTimeout` 不在 2000–60000 ms 范围内 |
| `InvalidScope` | Scope 为空、超过 4096 字符、超过 100 个值，或格式错误（非 `audience\|scopeValue`） |
| `MultipleAudienceNotSupported` | Scope 项引用了多个不同的 audience |

### 运行时错误

| 错误码 | 原因 |
|---|---|
| `AlibabaCloudCredentialError` | 阿里云默认凭证链未找到凭证 |
| `EIAMClientError` | EIAM SDK 客户端创建失败 |
| `GenerateOAuthTokenError` | EIAM `GenerateOauthToken` API 调用失败 |
| `InvalidTokenResponse` | EIAM 返回了不完整、为空或已过期的令牌 |

EIAM API 返回 HTTP 错误时，按以下方式分类：

| HTTP 状态码 | 错误类别 |
|---|---|
| 500+ | ServerError |
| 400–499 | ClientError |
| 其他 | HttpError |

## 开发验证

```bash
gofmt -w .
go vet ./...
go test -race -cover ./...
```

参见 [CONTRIBUTING.md](CONTRIBUTING.md)、[SECURITY.md](SECURITY.md) 和
[CHANGELOG.md](CHANGELOG.md)。

## 支持与反馈

- **邮箱**：cloudidaas@list.alibaba-inc.com
- **问题反馈**：如有疑问或建议，请在仓库提交 Issue。

## 许可证

本项目使用 Apache License 2.0，详见 [LICENSE](LICENSE) 和 [NOTICE](NOTICE)。
