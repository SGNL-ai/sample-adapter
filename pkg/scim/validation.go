// Copyright 2025 SGNL.ai, Inc.
package scim

import (
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(request *framework.Request[Config]) *framework.Error {
	if strings.HasPrefix(request.Address, "http://") {
		return &framework.Error{
			Message: "The provided HTTP protocol is not supported.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// SCIM server can use any of the Auth mechanisms
	if request.Auth == nil || (request.Auth.HTTPAuthorization == "" && request.Auth.Basic == nil) {
		return &framework.Error{
			Message: "SCIM auth is missing required credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth.Basic != nil && (request.Auth.Basic.Username == "" || request.Auth.Basic.Password == "") {
		return &framework.Error{
			Message: "One of username or password required for basic auth is empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Add checks for Ordered and MaxPageSize here, if any.
	// Depends on the SCIM server implementation hence excluded in the validation.

	return nil
}
