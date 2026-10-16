package alibabacloud

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/cloud-idaas/idaas-go-core-sdk/config"
	"github.com/cloud-idaas/idaas-go-core-sdk/enums"
	"github.com/cloud-idaas/idaas-go-core-sdk/plugin"
)

func TestCredentialProvider_GetName(t *testing.T) {
	t.Parallel()

	provider := NewCredentialProvider()
	if got := provider.GetName(); got != PluginName {
		t.Fatalf("GetName() = %q, want %q", got, PluginName)
	}
}

func TestAlibabaCloudPlugin_ZeroValuePreservesGetNameCompatibility(t *testing.T) {
	t.Parallel()

	provider := &AlibabaCloudPlugin{}
	if got := provider.GetName(); got != PluginName {
		t.Fatalf("GetName() = %q, want %q", got, PluginName)
	}
}

func TestPluginIsAutomaticallyRegistered(t *testing.T) {
	t.Parallel()

	registered, err := plugin.GetInstance().GetPlugin(PluginName)
	if err != nil {
		t.Fatalf("GetPlugin(%q) returned error: %v", PluginName, err)
	}
	if registered.GetName() != PluginName {
		t.Fatalf("registered plugin name = %q, want %q", registered.GetName(), PluginName)
	}
}

func TestReadmeConfigurationsUseGoPluginMethod(t *testing.T) {
	for _, path := range []string{"README.md", "README_zh.md"} {
		t.Run(path, func(t *testing.T) {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			section := strings.SplitN(string(content), "```json\n", 2)
			if len(section) != 2 {
				t.Fatal("README configuration example not found")
			}
			var cfg config.IDaaSClientConfig
			if err := json.Unmarshal([]byte(strings.SplitN(section[1], "\n```", 2)[0]), &cfg); err != nil {
				t.Fatal(err)
			}
			if err := config.NewConfigValidator().Validate(&cfg); err != nil {
				t.Fatalf("README configuration is invalid: %v", err)
			}
			if cfg.AuthnConfiguration.AuthnMethod != enums.TokenAuthnMethodPlugin ||
				cfg.AuthnConfiguration.PluginName != PluginName {
				t.Fatalf("README does not select this plugin: %#v", cfg.AuthnConfiguration)
			}
		})
	}
}
