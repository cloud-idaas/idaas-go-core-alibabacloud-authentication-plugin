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

	"github.com/alibabacloud-go/tea/dara"
	tea "github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
)

const (
	errCodeAlibabaCloudCredential = "AlibabaCloudCredentialError"
	errCodeEIAMClient             = "EIAMClientError"
	errCodeGenerateOAuthToken     = "GenerateOAuthTokenError"
	errCodeInvalidTokenResponse   = "InvalidTokenResponse"
)

func translateSDKError(err error) error {
	code := errCodeGenerateOAuthToken
	statusCode := 0

	var daraErr *dara.SDKError
	if stderrors.As(err, &daraErr) {
		if daraErr.Code != nil && *daraErr.Code != "" {
			code = *daraErr.Code
		}
		if daraErr.StatusCode != nil {
			statusCode = *daraErr.StatusCode
		}
	} else {
		var teaErr *tea.SDKError
		if stderrors.As(err, &teaErr) {
			if teaErr.Code != nil && *teaErr.Code != "" {
				code = *teaErr.Code
			}
			if teaErr.StatusCode != nil {
				statusCode = *teaErr.StatusCode
			}
		}
	}

	if statusCode >= 500 {
		return sdkerrors.NewServerError(code, "EIAM GenerateOauthToken request failed", statusCode, err)
	}
	if statusCode >= 400 {
		return sdkerrors.NewClientError(code, "EIAM GenerateOauthToken request failed", err)
	}
	return sdkerrors.NewHttpError(code, "EIAM GenerateOauthToken request failed", statusCode, err)
}
