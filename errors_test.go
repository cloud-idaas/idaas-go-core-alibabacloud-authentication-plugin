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
