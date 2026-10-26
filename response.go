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
	"strings"
	"time"

	eiam "github.com/alibabacloud-go/eiam-20211201/v2/client"
	"github.com/cloud-idaas/idaas-go-core-sdk/credential"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
)

func mapTokenResponse(response *eiam.GenerateOauthTokenResponse, requestedAt, receivedAt time.Time) (*credential.IDaaSTokenResponse, error) {
	if response == nil || response.Body == nil || response.Body.TokenResponse == nil {
		return nil, sdkerrors.NewCredentialError(
			errCodeInvalidTokenResponse,
			"EIAM returned an incomplete token response",
			nil,
		)
	}

	token := response.Body.TokenResponse
	if token.AccessToken == nil || strings.TrimSpace(*token.AccessToken) == "" {
		return nil, sdkerrors.NewCredentialError(
			errCodeInvalidTokenResponse,
			"EIAM returned an empty access token",
			nil,
		)
	}
	if token.ExpiresIn == nil || *token.ExpiresIn <= 0 {
		return nil, sdkerrors.NewCredentialError(
			errCodeInvalidTokenResponse,
			"EIAM returned an invalid token lifetime",
			nil,
		)
	}
	remainingSeconds := *token.ExpiresIn
	elapsed := receivedAt.Sub(requestedAt)
	if elapsed > 0 {
		elapsedSeconds := int64(elapsed / time.Second)
		if elapsed%time.Second != 0 {
			elapsedSeconds++
		}
		remainingSeconds -= elapsedSeconds
	}
	if token.ExpiresAt != nil {
		absoluteSeconds := int64(time.Unix(*token.ExpiresAt, 0).Sub(receivedAt) / time.Second)
		if absoluteSeconds < remainingSeconds {
			remainingSeconds = absoluteSeconds
		}
	}
	remainingSeconds--
	if remainingSeconds <= 0 {
		return nil, sdkerrors.NewCredentialError(
			errCodeInvalidTokenResponse,
			"EIAM returned an expired or unusable access token",
			nil,
		)
	}

	tokenType := "Bearer"
	if token.TokenType != nil && strings.TrimSpace(*token.TokenType) != "" {
		tokenType = strings.TrimSpace(*token.TokenType)
	}

	return &credential.IDaaSTokenResponse{
		AccessToken: strings.TrimSpace(*token.AccessToken),
		TokenType:   tokenType,
		ExpiresIn:   remainingSeconds,
		ReceivedAt:  receivedAt,
	}, nil
}
