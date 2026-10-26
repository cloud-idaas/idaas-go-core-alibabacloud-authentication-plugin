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

// Sample: PLUGIN authentication entry point.
//
// Obtains an OAuth2 token via the Alibaba Cloud EIAM OpenAPI
// GenerateOauthToken operation.
//
// Prerequisites:
//  1. Configure an Alibaba Cloud credential (environment variables
//     ALIBABA_CLOUD_ACCESS_KEY_ID / ALIBABA_CLOUD_ACCESS_KEY_SECRET)
//  2. Prepare an IDaaS config file (authnMethod=plugin,
//     pluginName=alibabacloudPluginCredentialProvider)
//
// Usage:
//
//	go run samples/plugin_authentication.go -config ~/.cloud_idaas/plugin-config.json
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cloud-idaas/idaas-go-core-sdk/config"
	"github.com/cloud-idaas/idaas-go-core-sdk/factory"

	// Import the plugin for its side effect to auto-register with PluginRegistry.
	_ "github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin"
)

func main() {
	configPath := flag.String("config", "", "path to IDaaS client config file")
	flag.Parse()

	reader := config.NewConfigReader()
	cfg, err := reader.LoadWithPriority(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	f := factory.GetInstance()
	if err = f.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize factory: %v", err)
	}

	provider, err := f.CreateCredentialProvider()
	if err != nil {
		log.Fatalf("Failed to create credential provider: %v", err)
	}

	authnMethod := ""
	if cfg.AuthnConfiguration != nil {
		authnMethod = string(cfg.AuthnConfiguration.AuthnMethod)
	}
	fmt.Printf("Getting access token (authnMethod=%s)...\n", authnMethod)

	cred, err := provider.GetCredential()
	if err != nil {
		log.Fatalf("Failed to get credential: %v", err)
	}

	fmt.Println("✅ Access Token obtained successfully!")
	fmt.Printf("   Token Type: %s\n", cred.GetTokenType())
	fmt.Printf("   Expires In: %d seconds\n", cred.GetExpiresIn())
	if scope := cred.GetScope(); scope != "" {
		fmt.Printf("   Scope:      %s\n", scope)
	}
}
