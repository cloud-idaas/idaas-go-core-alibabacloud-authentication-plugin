# 集成测试指南

## 前置条件

### 1. 配置阿里云凭证

设置环境变量（必须）：

```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID='你的AK'
export ALIBABA_CLOUD_ACCESS_KEY_SECRET='你的SK'
```

### 2. 准备 IDaaS 配置文件

复制模板并编辑：

```bash
cp config/integration-config.json.template config/integration-config.json
vim config/integration-config.json
```

配置文件示例：

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

**配置字段说明：**

| 字段 | 必填 | 说明 |
|------|------|------|
| idaasInstanceId | ✓ | IDaaS 实例 ID |
| clientId | ✓ | 应用 ID |
| issuer | ✓ | OAuth2 issuer 端点 |
| tokenEndpoint | ✓ | OAuth2 token 端点 |
| scope | ✓ | 权限范围（格式：`audience\|scopeValue`） |
| openApiEndpoint | ✓ | EIAM OpenAPI 端点 |

## 运行集成测试

### 方式 1：使用运行脚本（推荐）

```bash
./run_integration_test.sh
```

或指定配置文件：

```bash
./run_integration_test.sh config/integration-config.json
```

### 方式 2：直接运行 go test

```bash
go test -v -tags=integration -run TestIntegration
```

## 测试用例

### 1. TestIntegration_GetIDaaSCredential_RealAPI

调用真实的阿里云 EIAM API 获取 token，验证：
- 能够成功获取 access token
- token 类型正确
- 过期时间有效

### 2. TestIntegration_GetIDaaSCredential_InvalidCredentials

使用无效凭证调用 API，验证：
- 正确返回错误
- 错误信息清晰

## 常见问题

### Q: 提示"环境变量未设置"
A: 确保已设置 `ALIBABA_CLOUD_ACCESS_KEY_ID` 和 `ALIBABA_CLOUD_ACCESS_KEY_SECRET`

### Q: 提示"配置文件不存在"
A: 先复制模板：`cp config/integration-config.json.template config/integration-config.json`

### Q: 测试失败，返回认证错误
A: 检查：
1. AK/SK 是否正确
2. IDaaS 实例 ID 是否正确
3. 应用是否已授权对应的 scope
4. 网络是否可以访问 EIAM OpenAPI 端点

## 安全提示

⚠️ **重要**：
- 不要将 AK/SK 提交到代码仓库
- 不要将包含真实配置的文件提交到 Git
- 使用 `.gitignore` 忽略 `config/integration-config.json`
- 测试完成后及时清理环境变量

