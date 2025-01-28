// Copyright 2025 SGNL.ai, Inc.
package scim

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
)

// Client is a client that allows querying a SCIM SoR which
// contains JSON objects.
type Client interface {
	// GetPage returns a page of JSON objects from the datasource for the
	// requested entity.
	// Returns a (possibly empty) list of JSON objects, each object being
	// unmarshaled into a map by Golang's JSON unmarshaler.
	GetPage(ctx context.Context, request *Request) (*AdapterResponse, *framework.Error)
}

// Request is a request to a SCIM SoR.
type Request struct {
	// BaseURL is the Base URL of the datasource to query. For example, "my.scim.server.com".
	BaseURL string

	// AuthorizationHeader is the Authorization header sent to the SCIM SoR.
	AuthorizationHeader string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// Optional. If not set, return the first page for this entity.
	Cursor string

	// QueryParams contains the query parameters required to generate the URL for the datasource request
	QueryParams QueryParams

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
}

// AdapterResponse is a response returned by the adapter.
type AdapterResponse struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of items returned by the datasource.
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor string
}
