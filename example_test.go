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

package alibabacloud_test

import (
	"fmt"
	"log"

	"github.com/cloud-idaas/idaas-go-core-sdk/plugin"

	// Importing the plugin package for its side effect registers the
	// credential provider with the Core SDK plugin registry, which is
	// what the Core SDK relies on when authnConfiguration selects
	// pluginName "alibabacloudPluginCredentialProvider".
	alibabacloud "github.com/cloud-idaas/idaas-go-core-alibabacloud-authentication-plugin"
)

func Example() {
	registered, err := plugin.GetInstance().GetPlugin(alibabacloud.PluginName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(registered.GetName())
	// Output:
	// alibabacloudPluginCredentialProvider
}
