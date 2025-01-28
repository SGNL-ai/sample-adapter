// Copyright 2025 SGNL.ai, Inc.
package scim

import (
	"github.com/sgnl-ai/sample-adapter/pkg/config"
)

type QueryParams struct {
	// Filter allows to request a subset of resources via the "filter" query parameter containing a filter expression
	Filter string `json:"filter,omitempty"`

	// SortBy allows to sort the returned resources via the "sortBy" query parameter
	SortBy string `json:"sortBy,omitempty"`

	// Ascending allows to specify the sort order via the "sortOrder" query parameter
	Ascending *bool `json:"ascending,omitempty"`
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "queryParams": {
        "Users": {
            "filter": "userType eq \"Employee\" and (emails co \"sgnl.com\" or emails.value co \"sgnl.org\"",
            "sortBy": "userName",
            "ascending": true
        },
        "Groups": {
            "filter": "displayName eq \"SGNL\"",
            "sortBy": "displayName",
            "ascending": true
        }
    }
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// QueryParams is an map containing the query parameters for each entity associated with this
	// datasource. The key is the entity's external_name, and the value is the QueryParams.
	QueryParams map[string]QueryParams `json:"queryParams,omitempty"`
}
