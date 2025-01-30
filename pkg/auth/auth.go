// Copyright 2025 SGNL.ai, Inc.
package auth

import "encoding/base64"

// BasicAuthHeader returns HTTP Basic authentication credentials for the provided username and password.
func BasicAuthHeader(username, password string) string {
	authString := username + ":" + password

	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))
}
