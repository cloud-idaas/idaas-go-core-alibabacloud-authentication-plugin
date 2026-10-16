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
