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
	stderrors "errors"
	"testing"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	eiam "github.com/alibabacloud-go/eiam-20211201/v2/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/cloud-idaas/idaas-go-core-sdk/config"
	"github.com/cloud-idaas/idaas-go-core-sdk/constants"
	"github.com/cloud-idaas/idaas-go-core-sdk/enums"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
	"github.com/cloud-idaas/idaas-go-core-sdk/factory"
)

type fakeCloudCredential struct{}

func (fakeCloudCredential) GetAccessKeyId() (*string, error)     { return nil, nil }
func (fakeCloudCredential) GetAccessKeySecret() (*string, error) { return nil, nil }
func (fakeCloudCredential) GetSecurityToken() (*string, error)   { return nil, nil }
func (fakeCloudCredential) GetBearerToken() *string              { return nil }
func (fakeCloudCredential) GetType() *string                     { return nil }
func (fakeCloudCredential) GetCredential() (*credentials.CredentialModel, error) {
	return nil, nil
}

type fakeTokenClient struct {
	response *eiam.GenerateOauthTokenResponse
	err      error
	request  *eiam.GenerateOauthTokenRequest
	runtime  *dara.RuntimeOptions
}

func (f *fakeTokenClient) GenerateOauthTokenWithOptions(
	request *eiam.GenerateOauthTokenRequest,
	runtime *dara.RuntimeOptions,
) (*eiam.GenerateOauthTokenResponse, error) {
	f.request = request
	f.runtime = runtime
	return f.response, f.err
}

func TestCredentialProvider_GetIDaaSCredential_WhenSuccessful(t *testing.T) {
	client := &fakeTokenClient{
		response: tokenResponse(dara.String("token"), dara.Int64(3600), dara.String("Bearer")),
	}
	var capturedConfig *openapi.Config
	now := time.Unix(1_700_000_000, 0)
	provider := newCredentialProvider(
		func() (credentials.Credential, error) { return fakeCloudCredential{}, nil },
		func(cfg *openapi.Config) (eiamTokenClient, error) {
			capturedConfig = cfg
			return client, nil
		},
		func() time.Time { return now },
	)
	initializeFactory(t)

	got, err := provider.GetIDaaSCredential("api.example.com|write api.example.com|read")

	if err != nil {
		t.Fatalf("GetIDaaSCredential() returned error: %v", err)
	}
	if got.GetAccessToken() != "token" || got.GetExpiresIn() != 3599 {
		t.Fatalf("credential = %#v", got)
	}
	if capturedConfig == nil || dara.StringValue(capturedConfig.Endpoint) != "eiam.cn-hangzhou.aliyuncs.com" {
		t.Fatalf("client endpoint = %v", capturedConfig)
	}
	if dara.IntValue(capturedConfig.ConnectTimeout) != 2000 || dara.IntValue(capturedConfig.ReadTimeout) != 5000 {
		t.Fatalf("client timeouts were not propagated: %#v", capturedConfig)
	}
	if client.runtime == nil ||
		dara.IntValue(client.runtime.ConnectTimeout) != 2000 ||
		dara.IntValue(client.runtime.ReadTimeout) != 5000 {
		t.Fatalf("runtime timeouts were not propagated: %#v", client.runtime)
	}
	if dara.StringValue(client.request.InstanceId) != "idaas-instance" {
		t.Errorf("InstanceId = %q", dara.StringValue(client.request.InstanceId))
	}
	if dara.StringValue(client.request.ApplicationId) != "application" {
		t.Errorf("ApplicationId = %q", dara.StringValue(client.request.ApplicationId))
	}
	if dara.StringValue(client.request.Audience) != "api.example.com" {
		t.Errorf("Audience = %q", dara.StringValue(client.request.Audience))
	}
	if gotScopes := dara.StringSliceValue(client.request.ScopeValues); len(gotScopes) != 2 ||
		gotScopes[0] != "read" || gotScopes[1] != "write" {
		t.Errorf("ScopeValues = %v", gotScopes)
	}
}

func TestCredentialProvider_GetIDaaSCredential_AccountsForRequestDuration(t *testing.T) {
	client := &fakeTokenClient{
		response: tokenResponse(dara.String("token"), dara.Int64(10), nil),
	}
	requestedAt := time.Unix(1_700_000_000, 0)
	clockCalls := 0
	provider := newCredentialProvider(
		func() (credentials.Credential, error) { return fakeCloudCredential{}, nil },
		func(*openapi.Config) (eiamTokenClient, error) { return client, nil },
		func() time.Time {
			clockCalls++
			if clockCalls == 1 {
				return requestedAt
			}
			return requestedAt.Add(2500 * time.Millisecond)
		},
	)
	initializeFactory(t)

	got, err := provider.GetIDaaSCredential("api.example.com|read")
	if err != nil {
		t.Fatalf("GetIDaaSCredential() returned error: %v", err)
	}
	if got.GetExpiresIn() != 6 || clockCalls != 2 {
		t.Fatalf("expiresIn = %d, clock calls = %d; want 6, 2", got.GetExpiresIn(), clockCalls)
	}
}

func TestCredentialProvider_GetIDaaSCredential_WhenCredentialCreationFails(t *testing.T) {
	wantErr := stderrors.New("credential unavailable")
	provider := newCredentialProvider(
		func() (credentials.Credential, error) { return nil, wantErr },
		func(*openapi.Config) (eiamTokenClient, error) {
			t.Fatal("client factory must not be called")
			return nil, nil
		},
		time.Now,
	)
	initializeFactory(t)

	_, err := provider.GetIDaaSCredential("api.example.com|read")

	var credentialErr *sdkerrors.CredentialError
	if !stderrors.As(err, &credentialErr) {
		t.Fatalf("error = %T, want *errors.CredentialError", err)
	}
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("error does not wrap credential cause: %v", err)
	}
}

func TestCredentialProvider_GetIDaaSCredential_WhenClientCreationFails(t *testing.T) {
	wantErr := stderrors.New("client unavailable")
	provider := newCredentialProvider(
		func() (credentials.Credential, error) { return fakeCloudCredential{}, nil },
		func(*openapi.Config) (eiamTokenClient, error) { return nil, wantErr },
		time.Now,
	)
	initializeFactory(t)

	_, err := provider.GetIDaaSCredential("api.example.com|read")

	var clientErr *sdkerrors.ClientError
	if !stderrors.As(err, &clientErr) {
		t.Fatalf("error = %T, want *errors.ClientError", err)
	}
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("error does not wrap client cause: %v", err)
	}
}

func TestValidateHTTPTimeouts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		config   *config.HttpConfiguration
		wantCode string
	}{
		{name: "nil configuration"},
		{name: "defaults", config: &config.HttpConfiguration{}},
		{name: "minimum", config: &config.HttpConfiguration{ConnectTimeout: 2000, ReadTimeout: 2000}},
		{name: "maximum", config: &config.HttpConfiguration{ConnectTimeout: 60000, ReadTimeout: 60000}},
		{name: "connect too small", config: &config.HttpConfiguration{ConnectTimeout: 1999}, wantCode: constants.ConnectTimeoutNotValid},
		{name: "connect too large", config: &config.HttpConfiguration{ConnectTimeout: 60001}, wantCode: constants.ConnectTimeoutNotValid},
		{name: "read too small", config: &config.HttpConfiguration{ReadTimeout: -1}, wantCode: constants.ReadTimeoutNotValid},
		{name: "read too large", config: &config.HttpConfiguration{ReadTimeout: 60001}, wantCode: constants.ReadTimeoutNotValid},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validateHTTPTimeouts(test.config)
			if test.wantCode == "" {
				if err != nil {
					t.Fatalf("validateHTTPTimeouts() returned error: %v", err)
				}
				return
			}
			var configErr *sdkerrors.ConfigError
			if !stderrors.As(err, &configErr) {
				t.Fatalf("error = %T, want *errors.ConfigError", err)
			}
			if configErr.ErrorCode() != test.wantCode {
				t.Fatalf("error code = %q, want %q", configErr.ErrorCode(), test.wantCode)
			}
		})
	}
}

func initializeFactory(t *testing.T) {
	t.Helper()
	err := factory.GetInstance().Initialize(&config.IDaaSClientConfig{
		ClientId:        "application",
		InstanceId:      "idaas-instance",
		IssuerEndpoint:  "https://example.com/oauth2",
		TokenEndpoint:   "https://example.com/oauth2/token",
		OpenApiEndpoint: "eiam.cn-hangzhou.aliyuncs.com",
		Scope:           "api.example.com|read",
		AuthnConfiguration: &config.IdentityAuthenticationConfiguration{
			IdentityType: enums.AuthenticationIdentityClient,
			AuthnMethod:  enums.TokenAuthnMethodPlugin,
			PluginName:   PluginName,
		},
		HttpConfiguration: &config.HttpConfiguration{
			ConnectTimeout: 2000,
			ReadTimeout:    5000,
		},
	})
	if err != nil {
		t.Fatalf("Initialize() returned error: %v", err)
	}
}
