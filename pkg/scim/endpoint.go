// Copyright 2025 SGNL.ai, Inc.
package scim

import (
	"net/url"
	"strconv"
	"strings"
)

// GenerateURL returns a URL to fetch a given page of SCIM objects.
func GenerateURL(
	baseURL string,
	entityExternalID string,
	pageSize int64,
	startIndex string,
	queryParams QueryParams,
) string {
	escapedFilter := url.QueryEscape(queryParams.Filter)

	filterLen := len(escapedFilter)
	if filterLen > 0 {
		filterLen += 8 // len("&filter=") == 8
	}

	sortByLen := len(queryParams.SortBy)
	if sortByLen > 0 {
		sortByLen += 8 // len("&sortBy=") == 8
	}

	sortOrderLen := 0
	if queryParams.Ascending != nil {
		sortOrderLen += 21 // len("&sortOrder=") + max(len(descending), len(ascending)) == 21
	}

	// len(baseURL) + len("/") + len(entityExternalID) +
	// len("?count=") + len(strconv.FormatInt(pageSize, 10)) + len("&startIndex=") +
	// len(startIndex) +
	// filterLen + sortByLen + sortOrderLen ==

	// len(baseURL) + len(entityExternalID) +
	// len(strconv.FormatInt(pageSize, 10)) + len(strconv.FormatInt(startIndex, 10)) +
	// filterLen + sortByLen + sortOrderLen + 20
	var sb strings.Builder

	sb.Grow(
		len(baseURL) + len(entityExternalID) + len(strconv.FormatInt(pageSize, 10)) + len(startIndex) +
			filterLen + sortByLen + sortOrderLen + 20,
	)

	sb.WriteString(baseURL)
	sb.WriteString("/")
	sb.WriteString(entityExternalID)
	sb.WriteString("?startIndex=")
	sb.WriteString(startIndex)
	sb.WriteString("&count=")
	sb.WriteString(strconv.FormatInt(pageSize, 10))

	if queryParams.Filter != "" {
		sb.WriteString("&filter=")
		sb.WriteString(escapedFilter)
	}

	if queryParams.SortBy != "" {
		sb.WriteString("&sortBy=")
		sb.WriteString(queryParams.SortBy)
	}

	if queryParams.Ascending != nil {
		if *queryParams.Ascending {
			sb.WriteString("&sortOrder=ascending")
		} else {
			sb.WriteString("&sortOrder=descending")
		}
	}

	return sb.String()
}
