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

	"github.com/alibabacloud-go/tea/dara"
	tea "github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
)

func TestTranslateSDKError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      error
		wantCode   string
		wantServer bool
		wantHTTP   bool
	}{
		{
			name:     "Dara 4xx",
			input:    &dara.SDKError{Code: dara.String("Forbidden"), StatusCode: dara.Int(403)},
			wantCode: "Forbidden",
		},
		{
			name:       "Dara 5xx",
			input:      &dara.SDKError{Code: dara.String("InternalError"), StatusCode: dara.Int(503)},
			wantCode:   "InternalError",
			wantServer: true,
		},
		{
			name:     "Tea 4xx",
			input:    &tea.SDKError{Code: dara.String("BadRequest"), StatusCode: dara.Int(400)},
			wantCode: "BadRequest",
		},
		{
			name:     "generic error",
			input:    stderrors.New("network unavailable"),
			wantCode: errCodeGenerateOAuthToken,
			wantHTTP: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := translateSDKError(test.input)
			if !stderrors.Is(got, test.input) {
				t.Fatalf("translated error does not wrap input: %v", got)
			}
			if test.wantServer {
				var typed *sdkerrors.ServerError
				if !stderrors.As(got, &typed) {
					t.Fatalf("error = %T, want *errors.ServerError", got)
				}
				if typed.ErrorCode() != test.wantCode {
					t.Fatalf("error code = %q, want %q", typed.ErrorCode(), test.wantCode)
				}
				return
			}
			if test.wantHTTP {
				var typed *sdkerrors.HttpError
				if !stderrors.As(got, &typed) {
					t.Fatalf("error = %T, want *errors.HttpError", got)
				}
				if typed.ErrorCode() != test.wantCode {
					t.Fatalf("error code = %q, want %q", typed.ErrorCode(), test.wantCode)
				}
				return
			}

			var typed *sdkerrors.ClientError
			if !stderrors.As(got, &typed) {
				t.Fatalf("error = %T, want *errors.ClientError", got)
			}
			if typed.ErrorCode() != test.wantCode {
				t.Fatalf("error code = %q, want %q", typed.ErrorCode(), test.wantCode)
			}
		})
	}
}
