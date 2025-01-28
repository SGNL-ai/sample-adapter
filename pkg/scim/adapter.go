// Copyright 2025 SGNL.ai, Inc.
package scim

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapter-sgnl/pkg/auth"
	"github.com/sgnl-ai/adapter-sgnl/pkg/config"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from SCIM 2.0 datasources.
type Adapter struct {
	// Client provides access to the datasource.
	Client Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		Client: client,
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	if err := a.ValidateGetPageRequest(request); err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return a.RequestPageFromDatasource(ctx, request)
}

// RequestPageFromDatasource requests a page of objects from a SoR.
// It calls the SCIM SoR client internally to make the SoR request, parses the response,
// and handles any errors.
// It also handles parsing the current cursor and generating the next cursor.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context,
	request *framework.Request[Config],
) framework.Response {
	var commonConfig *config.CommonConfig
	if request.Config != nil {
		commonConfig = request.Config.CommonConfig
	}

	commonConfig = config.SetMissingCommonConfigDefaults(commonConfig)

	if !strings.HasPrefix(request.Address, "https://") {
		request.Address = "https://" + request.Address
	}

	var authorizationHeader string

	switch {
	case request.Auth.Basic != nil:
		authorizationHeader = auth.BasicAuthHeader(request.Auth.Basic.Username, request.Auth.Basic.Password)
	case request.Auth.HTTPAuthorization != "":
		authorizationHeader = request.Auth.HTTPAuthorization
	default:
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: "No valid credentials provided.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED,
			},
		)
	}

	req := &Request{
		BaseURL:               request.Address,
		AuthorizationHeader:   authorizationHeader,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		Cursor:                request.Cursor,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	if request.Config != nil && request.Config.QueryParams != nil {
		if entityQueryParams, found := request.Config.QueryParams[request.Entity.ExternalId]; found {
			req.QueryParams = entityQueryParams
		}
	}

	resp, err := a.Client.GetPage(ctx, req)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// Nested attributes are flattened and delimited by the delimiter specified.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,
		web.WithJSONPathAttributeNames(),
		web.WithLocalTimeZoneOffset(commonConfig.LocalTimeZoneOffset),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert SCIM response objects to JSON: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: resp.NextCursor,
	})
}
