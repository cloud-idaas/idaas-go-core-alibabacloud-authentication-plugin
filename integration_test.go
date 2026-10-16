//go:build integration
// +build integration

package alibabacloud

import (
	"os"
	"testing"

	"github.com/cloud-idaas/idaas-go-core-sdk/config"
	"github.com/cloud-idaas/idaas-go-core-sdk/factory"
)

// loadConfig 加载配置文件（使用 SDK 的 ConfigReader，会自动处理大小写归一化）
func loadConfig(t *testing.T) *config.IDaaSClientConfig {
	t.Helper()

	configPath := os.Getenv("IDAAS_CONFIG_FILE")
	if configPath == "" {
		configPath = "config/integration-config.json"
	}

	// 使用 SDK 的 ConfigReader 加载配置，会触发 UnmarshalJSON 的大小写归一化
	cfg, err := config.NewConfigReader().LoadFromFile(configPath)
	if err != nil {
		t.Skipf("配置文件加载失败: %s，跳过集成测试。设置 IDAAS_CONFIG_FILE 环境变量指定路径", err)
	}

	return cfg
}

// TestIntegration_GetIDaaSCredential_RealAPI 集成测试：调用真实的阿里云 EIAM API
func TestIntegration_GetIDaaSCredential_RealAPI(t *testing.T) {
	// 检查阿里云凭证
	if os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID") == "" || os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET") == "" {
		t.Skip("阿里云凭证未设置，跳过集成测试。请设置 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	}

	// 加载配置文件
	cfg := loadConfig(t)

	// 初始化 factory
	err := factory.GetInstance().Initialize(cfg)
	if err != nil {
		t.Fatalf("Initialize() 失败: %v", err)
	}

	// 创建 provider 并调用真实 API
	provider := NewCredentialProvider()

	t.Logf("正在调用 EIAM API 获取 token...")
	t.Logf("  InstanceId: %s", cfg.InstanceId)
	t.Logf("  ClientId: %s", cfg.ClientId)
	t.Logf("  Scope: %s", cfg.Scope)
	t.Logf("  Endpoint: %s", cfg.OpenApiEndpoint)

	cred, err := provider.GetIDaaSCredential(cfg.Scope)
	if err != nil {
		t.Fatalf("GetIDaaSCredential() 失败: %v", err)
	}

	// 验证返回的凭证
	if cred == nil {
		t.Fatal("返回的凭证为 nil")
	}

	accessToken := cred.GetAccessToken()
	if accessToken == "" {
		t.Error("AccessToken 为空")
	}

	expiresIn := cred.GetExpiresIn()
	if expiresIn <= 0 {
		t.Errorf("ExpiresIn 应该大于 0，实际值: %d", expiresIn)
	}

	tokenType := cred.GetTokenType()
	if tokenType == "" {
		t.Error("TokenType 为空")
	}

	t.Logf("✅ 成功获取 token:")
	t.Logf("   Token Type: %s", tokenType)
	t.Logf("   Expires In: %d 秒", expiresIn)
	t.Logf("   Access Token: [已隐藏，长度=%d]", len(accessToken))
}

// TestIntegration_GetIDaaSCredential_InvalidCredentials 集成测试：使用无效凭证
func TestIntegration_GetIDaaSCredential_InvalidCredentials(t *testing.T) {
	// 保存原始环境变量
	oldAK := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	oldSK := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")

	// 设置无效凭证
	os.Setenv("ALIBABA_CLOUD_ACCESS_KEY_ID", "invalid_key")
	os.Setenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET", "invalid_secret")
	defer func() {
		os.Setenv("ALIBABA_CLOUD_ACCESS_KEY_ID", oldAK)
		os.Setenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET", oldSK)
	}()

	// 加载配置文件
	cfg := loadConfig(t)

	// 初始化 factory
	err := factory.GetInstance().Initialize(cfg)
	if err != nil {
		t.Fatalf("Initialize() 失败: %v", err)
	}

	provider := NewCredentialProvider()
	_, err = provider.GetIDaaSCredential(cfg.Scope)

	// 使用无效凭证应该返回错误
	if err == nil {
		t.Error("使用无效凭证时应该返回错误，但返回了 nil")
	} else {
		t.Logf("✅ 正确返回错误: %v", err)
	}
}
