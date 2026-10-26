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
		{name: "single-label region", input: "eiam.c.aliyuncs.com", want: "eiam.c.aliyuncs.com"},
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
		{name: "trailing hyphen region", input: "eiam.cn-hangzhou-.aliyuncs.com", wantErr: true},
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
