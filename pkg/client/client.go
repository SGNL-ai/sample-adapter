// Copyright 2025 SGNL.ai, Inc.
package client

import (
	"net/http"
	"time"
)

type customTransport struct {
	userAgent string
	base      http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Set("User-Agent", t.userAgent)

	return t.base.RoundTrip(clonedReq)
}

func NewSGNLHttpClient(timeout time.Duration, userAgent string) *http.Client {
	// This is a default value if a SGNL 1P adapter does not define a custom user agent.
	if userAgent == "" {
		userAgent = "sgnl-adapter"
	}

	t := &customTransport{
		base:      http.DefaultTransport,
		userAgent: userAgent,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: t,
	}
}
