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
	"sort"
	"strings"

	"github.com/cloud-idaas/idaas-go-core-sdk/constants"
	sdkerrors "github.com/cloud-idaas/idaas-go-core-sdk/errors"
)

const (
	maxScopeLength = 4096
	maxScopeValues = 100
)

type audienceScope struct {
	audience    string
	scopeValues []string
}

func parseAudienceScope(raw string) (audienceScope, error) {
	if len(raw) > maxScopeLength {
		return audienceScope{}, invalidScopeError("scope exceeds the maximum length")
	}

	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return audienceScope{}, invalidScopeError("scope is empty")
	}
	if len(parts) > maxScopeValues {
		return audienceScope{}, invalidScopeError("scope contains too many values")
	}

	var audience string
	uniqueValues := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		if strings.Count(part, "|") != 1 {
			return audienceScope{}, invalidScopeError("scope item must use the audience|scopeValue format")
		}

		itemAudience, value, _ := strings.Cut(part, "|")
		if itemAudience == "" || value == "" {
			return audienceScope{}, invalidScopeError("scope audience and value must not be empty")
		}
		if audience == "" {
			audience = itemAudience
		} else if audience != itemAudience {
			return audienceScope{}, sdkerrors.NewConfigError(
				constants.MultipleAudienceNotSupported,
				"multiple audiences are not supported",
				nil,
			)
		}
		uniqueValues[value] = struct{}{}
	}

	values := make([]string, 0, len(uniqueValues))
	for value := range uniqueValues {
		values = append(values, value)
	}
	sort.Strings(values)

	return audienceScope{audience: audience, scopeValues: values}, nil
}

func invalidScopeError(message string) error {
	return sdkerrors.NewConfigError(constants.InvalidScope, message, nil)
}
