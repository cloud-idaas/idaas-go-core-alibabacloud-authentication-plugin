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
	"fmt"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	eiam "github.com/alibabacloud-go/eiam-20211201/v2/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/cloud-idaas/idaas-go-core-sdk/config"
	"github.com/cloud-idaas/idaas-go-core-sdk/constants"
	"github.com/cloud-idaas/idaas-go-core-sdk/credential"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
	"github.com/cloud-idaas/idaas-go-core-sdk/factory"
)

type eiamTokenClient interface {
	GenerateOauthTokenWithOptions(
		request *eiam.GenerateOauthTokenRequest,
		runtime *dara.RuntimeOptions,
	) (*eiam.GenerateOauthTokenResponse, error)
}

type credentialFactory func() (credentials.Credential, error)
type clientFactory func(*openapi.Config) (eiamTokenClient, error)

const (
	minHTTPTimeoutMillis = 2000
	maxHTTPTimeoutMillis = 60000
)

// CredentialProvider obtains IDaaS machine credentials through Alibaba Cloud
// EIAM. Token caching and refresh coordination are handled by the Core SDK.
type CredentialProvider struct {
	newCredential credentialFactory
	newClient     clientFactory
	now           func() time.Time
}

// AlibabaCloudPlugin is retained as a compatibility alias.
type AlibabaCloudPlugin = CredentialProvider

// NewCredentialProvider creates a provider backed by the Alibaba Cloud default
// credential chain and the official EIAM SDK.
func NewCredentialProvider() *CredentialProvider {
	return newCredentialProvider(
		defaultCredentialFactory,
		defaultClientFactory,
		time.Now,
	)
}

func defaultCredentialFactory() (credentials.Credential, error) {
	return credentials.NewCredential(nil)
}

func defaultClientFactory(config *openapi.Config) (eiamTokenClient, error) {
	return eiam.NewClient(config)
}

func newCredentialProvider(
	newCredential credentialFactory,
	newClient clientFactory,
	now func() time.Time,
) *CredentialProvider {
	return &CredentialProvider{
		newCredential: newCredential,
		newClient:     newClient,
		now:           now,
	}
}

// GetName returns the stable name used to select this plugin in IDaaS
// configuration.
func (p *CredentialProvider) GetName() string {
	return PluginName
}

// GetIDaaSCredential obtains an OAuth2 access token from Alibaba Cloud EIAM.
// Scope must contain one or more "audience|scopeValue" items sharing the same
// audience. The method never logs credential material; the Core SDK owns token
// caching and refresh behavior.
func (p *CredentialProvider) GetIDaaSCredential(scope string) (credential.IDaaSCredential, error) {
	cfg := factory.GetInstance().GetConfig()
	if cfg == nil {
		return nil, sdkerrors.NewConfigError(
			constants.IDaaSCredentialProviderFactoryNotInit,
			"IDaaS credential provider factory is not initialized",
			nil,
		)
	}
	if strings.TrimSpace(cfg.InstanceId) == "" {
		return nil, sdkerrors.NewConfigError(constants.IDaaSInstanceIdNotFound, "IDaaS instance ID is required", nil)
	}
	if strings.TrimSpace(cfg.ClientId) == "" {
		return nil, sdkerrors.NewConfigError(constants.ClientIdNotFound, "client ID is required", nil)
	}

	parsedScope, err := parseAudienceScope(scope)
	if err != nil {
		return nil, err
	}
	endpoint, err := parseAndValidateEndpoint(cfg.OpenApiEndpoint)
	if err != nil {
		return nil, err
	}
	if err := validateHTTPTimeouts(cfg.HttpConfiguration); err != nil {
		return nil, err
	}

	newCredential := p.newCredential
	if newCredential == nil {
		newCredential = defaultCredentialFactory
	}
	newClient := p.newClient
	if newClient == nil {
		newClient = defaultClientFactory
	}
	now := p.now
	if now == nil {
		now = time.Now
	}

	cloudCredential, err := newCredential()
	if err != nil {
		return nil, sdkerrors.NewCredentialError(
			errCodeAlibabaCloudCredential,
			"failed to obtain an Alibaba Cloud credential",
			err,
		)
	}

	clientConfig := &openapi.Config{
		Credential: cloudCredential,
		Endpoint:   dara.String(endpoint),
		Protocol:   dara.String("HTTPS"),
	}
	runtime := &dara.RuntimeOptions{}
	if cfg.HttpConfiguration != nil {
		if cfg.HttpConfiguration.ConnectTimeout > 0 {
			clientConfig.ConnectTimeout = dara.Int(cfg.HttpConfiguration.ConnectTimeout)
			runtime.ConnectTimeout = dara.Int(cfg.HttpConfiguration.ConnectTimeout)
		}
		if cfg.HttpConfiguration.ReadTimeout > 0 {
			clientConfig.ReadTimeout = dara.Int(cfg.HttpConfiguration.ReadTimeout)
			runtime.ReadTimeout = dara.Int(cfg.HttpConfiguration.ReadTimeout)
		}
	}

	client, err := newClient(clientConfig)
	if err != nil {
		return nil, sdkerrors.NewClientError(errCodeEIAMClient, "failed to create the EIAM client", err)
	}

	request := &eiam.GenerateOauthTokenRequest{
		InstanceId:    dara.String(cfg.InstanceId),
		ApplicationId: dara.String(cfg.ClientId),
		Audience:      dara.String(parsedScope.audience),
		ScopeValues:   dara.StringSlice(parsedScope.scopeValues),
	}
	requestedAt := now()
	response, err := client.GenerateOauthTokenWithOptions(request, runtime)
	if err != nil {
		return nil, translateSDKError(err)
	}

	return mapTokenResponse(response, requestedAt, now())
}

func validateHTTPTimeouts(httpConfig *config.HttpConfiguration) error {
	if httpConfig == nil {
		return nil
	}
	if err := validateOptionalTimeout(httpConfig.ConnectTimeout, constants.ConnectTimeoutNotValid, "connectTimeout"); err != nil {
		return err
	}
	return validateOptionalTimeout(httpConfig.ReadTimeout, constants.ReadTimeoutNotValid, "readTimeout")
}

func validateOptionalTimeout(value int, code, field string) error {
	if value == 0 {
		return nil
	}
	if value < minHTTPTimeoutMillis || value > maxHTTPTimeoutMillis {
		return sdkerrors.NewConfigError(
			code,
			fmt.Sprintf("%s must be between %d and %d milliseconds", field, minHTTPTimeoutMillis, maxHTTPTimeoutMillis),
			nil,
		)
	}
	return nil
}
