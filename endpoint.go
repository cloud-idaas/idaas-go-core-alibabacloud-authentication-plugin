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
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/cloud-idaas/idaas-go-core-sdk/constants"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
)

var eiamEndpointPattern = regexp.MustCompile("^eiam(?:-vpc)?\\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.aliyuncs\\.com$")

func parseAndValidateEndpoint(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", sdkerrors.NewConfigError(
			constants.OpenApiEndpointNotFound,
			"Alibaba Cloud EIAM OpenAPI endpoint is required",
			nil,
		)
	}

	endpointURL := raw
	if !strings.Contains(raw, "://") {
		endpointURL = "https://" + raw
	}
	u, err := url.Parse(endpointURL)
	if err != nil {
		return "", invalidEndpointError("invalid EIAM OpenAPI endpoint", err)
	}
	if !strings.EqualFold(u.Scheme, "https") {
		return "", invalidEndpointError("EIAM OpenAPI endpoint must use HTTPS", nil)
	}
	if u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return "", invalidEndpointError("EIAM OpenAPI endpoint must contain only a host", nil)
	}

	hostname := strings.ToLower(u.Hostname())
	if hostname == "" || hostname == "localhost" || net.ParseIP(hostname) != nil {
		return "", invalidEndpointError("EIAM OpenAPI endpoint host is not allowed", nil)
	}
	if port := u.Port(); port != "" && port != "443" {
		return "", invalidEndpointError("EIAM OpenAPI endpoint only supports port 443", nil)
	}
	if !eiamEndpointPattern.MatchString(hostname) {
		return "", invalidEndpointError("EIAM OpenAPI endpoint is not an official Alibaba Cloud endpoint", nil)
	}

	if u.Port() == "443" {
		return net.JoinHostPort(hostname, "443"), nil
	}
	return hostname, nil
}

func invalidEndpointError(message string, cause error) error {
	return sdkerrors.NewConfigError(constants.OpenApiEndpointNotFound, message, cause)
}
