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
	"testing"
	"time"

	eiam "github.com/alibabacloud-go/eiam-20211201/v2/client"
	"github.com/alibabacloud-go/tea/dara"
)

func TestMapTokenResponse(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	tests := []struct {
		name         string
		response     *eiam.GenerateOauthTokenResponse
		wantToken    string
		wantType     string
		wantExpiry   int64
		wantErr      bool
		requestDelay time.Duration
		expiresAt    *int64
	}{
		{name: "nil response", wantErr: true},
		{name: "nil body", response: &eiam.GenerateOauthTokenResponse{}, wantErr: true},
		{
			name:     "nil token response",
			response: &eiam.GenerateOauthTokenResponse{Body: &eiam.GenerateOauthTokenResponseBody{}},
			wantErr:  true,
		},
		{
			name:     "empty access token",
			response: tokenResponse(dara.String(" "), dara.Int64(3600), dara.String("Bearer")),
			wantErr:  true,
		},
		{
			name:     "missing expiry",
			response: tokenResponse(dara.String("token"), nil, dara.String("Bearer")),
			wantErr:  true,
		},
		{
			name:     "non-positive expiry",
			response: tokenResponse(dara.String("token"), dara.Int64(0), dara.String("Bearer")),
			wantErr:  true,
		},
		{
			name:       "defaults token type",
			response:   tokenResponse(dara.String(" token "), dara.Int64(3600), nil),
			wantToken:  "token",
			wantType:   "Bearer",
			wantExpiry: 3599,
		},
		{
			name:       "maps complete response",
			response:   tokenResponse(dara.String("token"), dara.Int64(7200), dara.String("DPoP")),
			wantToken:  "token",
			wantType:   "DPoP",
			wantExpiry: 7199,
		},
		{
			name:         "accounts for request latency",
			response:     tokenResponse(dara.String("token"), dara.Int64(3600), nil),
			requestDelay: 2500 * time.Millisecond,
			wantToken:    "token",
			wantType:     "Bearer",
			wantExpiry:   3596,
		},
		{
			name:       "uses earlier absolute expiry",
			response:   tokenResponse(dara.String("token"), dara.Int64(3600), nil),
			expiresAt:  dara.Int64(now.Add(30 * time.Second).Unix()),
			wantToken:  "token",
			wantType:   "Bearer",
			wantExpiry: 29,
		},
		{
			name:       "later absolute expiry does not extend lifetime",
			response:   tokenResponse(dara.String("token"), dara.Int64(30), nil),
			expiresAt:  dara.Int64(now.Add(time.Hour).Unix()),
			wantToken:  "token",
			wantType:   "Bearer",
			wantExpiry: 29,
		},
		{
			name:      "expired absolute timestamp",
			response:  tokenResponse(dara.String("token"), dara.Int64(3600), nil),
			expiresAt: dara.Int64(now.Add(-time.Second).Unix()),
			wantErr:   true,
		},
		{
			name:     "lifetime too short for cache",
			response: tokenResponse(dara.String("token"), dara.Int64(1), nil),
			wantErr:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.expiresAt != nil {
				test.response.Body.TokenResponse.ExpiresAt = test.expiresAt
			}
			got, err := mapTokenResponse(test.response, now.Add(-test.requestDelay), now)
			if test.wantErr {
				if err == nil {
					t.Fatal("mapTokenResponse() returned nil error")
				}
				return
			}
			if err != nil {
				t.Fatalf("mapTokenResponse() returned error: %v", err)
			}
			if got.AccessToken != test.wantToken || got.TokenType != test.wantType || got.ExpiresIn != test.wantExpiry {
				t.Fatalf("mapped token = %#v", got)
			}
			if !got.ReceivedAt.Equal(now) {
				t.Fatalf("ReceivedAt = %v, want %v", got.ReceivedAt, now)
			}
		})
	}
}

func tokenResponse(
	accessToken *string,
	expiresIn *int64,
	tokenType *string,
) *eiam.GenerateOauthTokenResponse {
	return &eiam.GenerateOauthTokenResponse{
		Body: &eiam.GenerateOauthTokenResponseBody{
			TokenResponse: &eiam.GenerateOauthTokenResponseBodyTokenResponse{
				AccessToken: accessToken,
				ExpiresIn:   expiresIn,
				TokenType:   tokenType,
			},
		},
	}
}
