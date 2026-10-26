# Integration Test Guide

> 本文档的中文版见 [INTEGRATION_TEST_zh.md](INTEGRATION_TEST_zh.md)。

Integration tests exercise the plugin against the real Alibaba Cloud EIAM
`GenerateOauthToken` API. They are guarded by the `integration` build tag and
are skipped unless credentials and a configuration file are provided.

## Prerequisites

### 1. Alibaba Cloud credentials

Export the credentials used by the Alibaba Cloud default credential chain
(required):

```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID='<your AK>'
export ALIBABA_CLOUD_ACCESS_KEY_SECRET='<your SK>'
```

### 2. IDaaS configuration file

Copy the template and edit it:

```bash
cp config/integration-config.json.template config/integration-config.json
vim config/integration-config.json
```

Example configuration:

```json
{
  "idaasInstanceId": "idaas_xxx",
  "clientId": "app_xxx",
  "issuer": "https://your-instance.aliyunidaas.com/api/v2/iauths_system/oauth2",
  "tokenEndpoint": "https://your-instance.aliyunidaas.com/api/v2/iauths_system/oauth2/token",
  "scope": "api.example.com|read api.example.com|write",
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

| Field | Required | Description |
|------|------|------|
| idaasInstanceId | ✓ | IDaaS instance ID |
| clientId | ✓ | Application ID |
| issuer | ✓ | OAuth2 issuer endpoint |
| tokenEndpoint | ✓ | OAuth2 token endpoint |
| scope | ✓ | Grant scope (format: `audience\|scopeValue`) |
| openApiEndpoint | ✓ | EIAM OpenAPI endpoint (a bare host is accepted and normalized to `https://` at runtime) |
| authnMethod | ✓ | Authentication method, case-insensitive. The IDaaS configuration API returns the uppercase form `PLUGIN`; the canonical Core SDK enum value is lowercase `plugin` |
| pluginName | ✓ | Must be `alibabacloudPluginCredentialProvider`; with any other value the Core SDK cannot resolve this plugin |
| httpConfiguration.connectTimeout / readTimeout | ✗ | Optional, in milliseconds; valid range 2000-60000, `0` inherits the default |

## Running the tests

### Option 1: helper script (recommended)

```bash
./run_integration_test.sh
```

Or with an explicit config path:

```bash
./run_integration_test.sh config/integration-config.json
```

### Option 2: plain `go test`

```bash
go test -v -tags=integration -run TestIntegration
```

## Test cases

### 1. TestIntegration_GetIDaaSCredential_RealAPI

Calls the real EIAM API to obtain a token and verifies:
- an access token is returned
- the token type is correct
- the expiry is valid

### 2. TestIntegration_GetIDaaSCredential_InvalidCredentials

Calls the API with invalid credentials and verifies:
- an error is returned
- the error message is clear

## Troubleshooting

### "Alibaba Cloud credentials are not set"
Set `ALIBABA_CLOUD_ACCESS_KEY_ID` and `ALIBABA_CLOUD_ACCESS_KEY_SECRET`.

### "config file not found"
Copy the template first:
`cp config/integration-config.json.template config/integration-config.json`

### Authentication errors in the test run
Check:
1. the AK/SK pair is valid
2. the IDaaS instance ID is correct
3. the application is granted the requested scopes
4. the network can reach the EIAM OpenAPI endpoint

## Security notes

⚠️ **Important**:
- Never commit AK/SK to the repository
- Never commit a config file containing real values; `config/integration-config.json`
  is gitignored on purpose
- Clear the credential environment variables after the run
