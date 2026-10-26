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

// Package alibabacloud provides an Alibaba Cloud EIAM credential-provider
// plugin for the IDaaS Go Core SDK.
//
// Import the package for its side effect to register the provider:
//
//	import _ "github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin"
package alibabacloud

import "github.com/cloud-idaas/idaas-go-core-sdk/plugin"

// PluginName is the registry name used in IDaaS Core SDK configuration.
const PluginName = "alibabacloudPluginCredentialProvider"

func init() {
	if err := plugin.GetInstance().Register(NewCredentialProvider()); err != nil {
		panic("alibabacloud credential provider registration failed: " + err.Error())
	}
}
