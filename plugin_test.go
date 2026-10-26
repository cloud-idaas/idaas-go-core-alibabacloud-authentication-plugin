// Copyright 2026 Alibaba Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
			// README examples intentionally use a bare openApiEndpoint host.
			// In production, factory.Initialize normalizes endpoints (a bare
			// host defaults to https://) before validating, so the config
			// package docs require callers to run NormalizeEndpoints() before
			// Validate as well. Mirror that Load -> Normalize -> Validate
			// pipeline here.
			cfg.NormalizeEndpoints()
			if err := config.NewConfigValidator().Validate(&cfg); err != nil {
				t.Fatalf("README configuration is invalid: %v", err)
			}
			if _, err := parseAndValidateEndpoint(cfg.OpenApiEndpoint); err != nil {
				t.Fatalf("README endpoint rejected by plugin: %v", err)
			}
			if cfg.AuthnConfiguration.AuthnMethod != enums.TokenAuthnMethodPlugin ||
				cfg.AuthnConfiguration.PluginName != PluginName {
				t.Fatalf("README does not select this plugin: %#v", cfg.AuthnConfiguration)
			}
		})
	}
}
