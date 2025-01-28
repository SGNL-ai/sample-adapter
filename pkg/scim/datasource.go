// Copyright 2025 SGNL.ai, Inc.
package scim

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapter-sgnl/pkg/errors"
)

// Datasource directly implements a Client interface to allow querying
// an external datasource.
type Datasource struct {
	Client *http.Client
}

type Response struct {
	Resources    []map[string]any `json:"Resources"`
	TotalResults int64            `json:"totalResults"`
	StartIndex   int64            `json:"startIndex"`
	ItemsPerPage int64            `json:"itemsPerPage"`
}

// NewClient instantiates and returns a new SCIM Client used to query the SCIM datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

// GetPage makes a request to the SCIM SoR to get a page of JSON objects. If a response is received,
// regardless of status code, a Response object is returned with the response body and the status code.
// If the request fails, an appropriate framework.Error is returned.
func (d *Datasource) GetPage(ctx context.Context, request *Request) (*AdapterResponse, *framework.Error) {
	cursor := "1"
	if request.Cursor != "" {
		cursor = request.Cursor
	}

	url := GenerateURL(
		request.BaseURL,
		request.EntityExternalID,
		request.PageSize,
		cursor,
		request.QueryParams,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: "Failed to create HTTP request to datasource.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)
	req.Header.Add("Accept", "application/scim+json")
	req.Header.Add("Authorization", request.AuthorizationHeader)

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute SCIM request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	defer res.Body.Close()

	response := &AdapterResponse{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	if res.StatusCode != http.StatusOK {
		return response, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: "Failed to read response body.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	objects, nextCursor, frameworkErr := ParseResponse(body, request.PageSize)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.Objects = objects
	response.NextCursor = nextCursor

	return response, nil
}

func ParseResponse(body []byte, pageSize int64) (objects []map[string]any, nextCursor string, err *framework.Error) {
	var scimResponse *Response

	if unmarshalErr := json.Unmarshal(body, &scimResponse); unmarshalErr != nil {
		return nil, "", &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if scimResponse.ItemsPerPage > pageSize {
		return nil, "", &framework.Error{
			Message: fmt.Sprintf("SCIM SoR returned more than the requested page size: %v.", scimResponse.ItemsPerPage),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	nextCursor = ""

	nextStartIndex := scimResponse.StartIndex + scimResponse.ItemsPerPage
	if nextStartIndex <= scimResponse.TotalResults {
		nextCursor = strconv.FormatInt(nextStartIndex, 10)
	}

	return scimResponse.Resources, nextCursor, nil
}
