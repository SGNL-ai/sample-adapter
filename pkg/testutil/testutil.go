// Copyright 2025 SGNL.ai, Inc.
package testutil

func GenPtr[T any](v T) *T {
	return &v
}
