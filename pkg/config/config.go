// Copyright 2025 SGNL.ai, Inc.
package config

var (
	DefaultRequestTimeout = 10 // 10 seconds
)

// CommonConfig is a collection of configuration common to all adapters.
type CommonConfig struct {
	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds *int `json:"requestTimeoutSeconds" validate:"omitempty,gt=0,lte=600"` // 10 minutes

	// LocalTimeZoneOffset is the default local timezone offset that should be used for
	// parsing date-time attributes lacking any time zone info. This should be set to the
	// number of seconds east of UTC. If this is set to 0 or not set, this will default to UTC.
	// Allowed offset is -12 hours to 14 hours, in seconds.
	LocalTimeZoneOffset int `json:"localTimeZoneOffset,omitempty" validate:"omitempty,gte=-43200,lte=50400"`
}

// SetMissingCommonConfigDefaults sets default values for any missing common configuration values.
// If the provided CommonConfig is nil, a new CommonConfig will be created.
func SetMissingCommonConfigDefaults(c *CommonConfig) *CommonConfig {
	if c == nil {
		c = &CommonConfig{}
	}

	// Set default request timeout.
	if c.RequestTimeoutSeconds == nil {
		c.RequestTimeoutSeconds = &DefaultRequestTimeout
	}

	return c
}
