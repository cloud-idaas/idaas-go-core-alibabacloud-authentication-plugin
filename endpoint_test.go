package alibabacloud

import "testing"

func TestParseAndValidateEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "public host", input: "eiam.cn-hangzhou.aliyuncs.com", want: "eiam.cn-hangzhou.aliyuncs.com"},
		{name: "HTTPS URL", input: " https://EIAM.cn-hangzhou.aliyuncs.com ", want: "eiam.cn-hangzhou.aliyuncs.com"},
		{name: "VPC host", input: "eiam-vpc.cn-shanghai.aliyuncs.com", want: "eiam-vpc.cn-shanghai.aliyuncs.com"},
		{name: "explicit HTTPS port", input: "https://eiam.cn-hangzhou.aliyuncs.com:443", want: "eiam.cn-hangzhou.aliyuncs.com:443"},
		{name: "empty", input: "", wantErr: true},
		{name: "HTTP", input: "http://eiam.cn-hangzhou.aliyuncs.com", wantErr: true},
		{name: "non EIAM domain", input: "https://example.com", wantErr: true},
		{name: "lookalike domain", input: "https://eiam.cn-hangzhou.aliyuncs.com.example.com", wantErr: true},
		{name: "IP literal", input: "https://127.0.0.1", wantErr: true},
		{name: "metadata IP", input: "https://169.254.169.254", wantErr: true},
		{name: "userinfo", input: "https://user@eiam.cn-hangzhou.aliyuncs.com", wantErr: true},
		{name: "path", input: "https://eiam.cn-hangzhou.aliyuncs.com/token", wantErr: true},
		{name: "root path", input: "https://eiam.cn-hangzhou.aliyuncs.com/", wantErr: true},
		{name: "query", input: "https://eiam.cn-hangzhou.aliyuncs.com?x=1", wantErr: true},
		{name: "fragment", input: "https://eiam.cn-hangzhou.aliyuncs.com#x", wantErr: true},
		{name: "unsupported port", input: "https://eiam.cn-hangzhou.aliyuncs.com:8443", wantErr: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseAndValidateEndpoint(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatalf("parseAndValidateEndpoint(%q) returned nil error", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseAndValidateEndpoint(%q) returned error: %v", test.input, err)
			}
			if got != test.want {
				t.Fatalf("parseAndValidateEndpoint(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func FuzzParseAndValidateEndpoint(f *testing.F) {
	for _, seed := range []string{
		"eiam.cn-hangzhou.aliyuncs.com",
		"https://eiam-vpc.cn-shanghai.aliyuncs.com",
		"http://127.0.0.1",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_, _ = parseAndValidateEndpoint(input)
	})
}
