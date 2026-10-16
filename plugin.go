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
