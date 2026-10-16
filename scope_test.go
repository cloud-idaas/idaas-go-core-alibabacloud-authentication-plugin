package alibabacloud

import (
	"strings"
	"testing"
)

func TestParseAudienceScope(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        audienceScope
		wantErr     bool
		wantErrCode string
	}{
		{
			name:  "single value",
			input: "api.example.com|read:file",
			want:  audienceScope{audience: "api.example.com", scopeValues: []string{"read:file"}},
		},
		{
			name:  "sorts and deduplicates values",
			input: "api.example.com|write api.example.com|read api.example.com|write",
			want:  audienceScope{audience: "api.example.com", scopeValues: []string{"read", "write"}},
		},
		{name: "empty", input: "", wantErr: true, wantErrCode: "InvalidScope"},
		{name: "whitespace only", input: " \t\n ", wantErr: true, wantErrCode: "InvalidScope"},
		{name: "missing separator", input: "read", wantErr: true, wantErrCode: "InvalidScope"},
		{name: "too many separators", input: "api|read|write", wantErr: true, wantErrCode: "InvalidScope"},
		{name: "empty audience", input: "|read", wantErr: true, wantErrCode: "InvalidScope"},
		{name: "empty value", input: "api|", wantErr: true, wantErrCode: "InvalidScope"},
		{
			name:        "multiple audiences",
			input:       "api-one|read api-two|write",
			wantErr:     true,
			wantErrCode: "MultipleAudienceNotSupported",
		},
		{
			name:        "too long",
			input:       strings.Repeat("a", maxScopeLength+1),
			wantErr:     true,
			wantErrCode: "InvalidScope",
		},
		{
			name:        "too many values",
			input:       strings.Repeat("api|read ", maxScopeValues+1),
			wantErr:     true,
			wantErrCode: "InvalidScope",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseAudienceScope(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatal("parseAudienceScope() returned nil error")
				}
				if coded, ok := err.(interface{ ErrorCode() string }); !ok || coded.ErrorCode() != test.wantErrCode {
					t.Fatalf("error code = %v, want %q", err, test.wantErrCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseAudienceScope() returned error: %v", err)
			}
			if got.audience != test.want.audience {
				t.Errorf("audience = %q, want %q", got.audience, test.want.audience)
			}
			if strings.Join(got.scopeValues, ",") != strings.Join(test.want.scopeValues, ",") {
				t.Errorf("scopeValues = %v, want %v", got.scopeValues, test.want.scopeValues)
			}
		})
	}
}

func FuzzParseAudienceScope(f *testing.F) {
	for _, seed := range []string{"api|read", "", "api|read|write", "api|read api|write"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_, _ = parseAudienceScope(input)
	})
}
