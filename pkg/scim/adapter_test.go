// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package scim_test

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-sgnl/pkg/config"
	"github.com/sgnl-ai/adapter-sgnl/pkg/scim"
	"github.com/sgnl-ai/adapter-sgnl/pkg/testutil"
)

var (
	testUsername = "mockuser"
	testPassword = "mockpassword"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	baseURL := server.URL
	adapter := scim.NewAdapter(&scim.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[scim.Config]
		wantResponse framework.Response
	}{
		"valid_user_request_without_common_config": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: baseURL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "userName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "emails",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "type",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "primary",
									Type:       framework.AttributeTypeBool,
								},
							},
						},
					},
				},

				PageSize: 2,
				Cursor:   "5",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "2819c223-7f76-453a-919d-413861904000",
							"userName": "bjensen@example.com",
							"emails": []framework.Object{
								{
									"value":   "bjensen@example.com",
									"type":    "work",
									"primary": true,
								},
								{
									"value": "babs@jensen.org",
									"type":  "home",
								},
							},
						},
					},
					NextCursor: "",
				},
			},
		},
		"valid_group_request_without_common_config": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: baseURL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimGroup,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.meta.created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "members",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "type",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},

				PageSize: 2,
				Cursor:   "5",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":             "e9e30dba-f08f-4109-8486-d5c6a331660a",
							"displayName":    "Tour Guides",
							"$.meta.created": time.Date(2010, 1, 23, 4, 56, 22, 0, time.UTC),
							"members": []framework.Object{
								{
									"value": "2819c223-7f76-453a-919d-413861904646",
									"type":  "User",
								},
								{
									"value": "6c5bb468-14b2-4183-baf2-06d523e03bd3",
									"type":  "Group",
								},
							},
						},
					},
					NextCursor: "",
				},
			},
		},
		"group_membership_from_valid_user_request": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: baseURL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "userName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "groups",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "display",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},

				PageSize: 2,
				Cursor:   "5", // startIndex=5 contains Full Enterprise User Extension Representation
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "2819c223-7f76-453a-919d-413861904000",
							"userName": "bjensen@example.com",
							"groups": []framework.Object{
								{
									"value":   "e9e30dba-f08f-4109-8486-d5c6a331660a",
									"display": "Tour Guides",
								},
								{
									"value":   "fc348aa8-3835-40eb-a20b-c726e15c55b5",
									"display": "Employees",
								},
								{
									"value":   "71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
									"display": "US Employees",
								},
							},
						},
					},
					NextCursor: "",
				},
			},
		},
		"invalid_request_missing_auth": {
			request: &framework.Request[scim.Config]{
				Address: "example.com",
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "SCIM auth is missing required credentials.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		// This test ensures that if the SCIM SoR returns a non successful status code, we return an
		// appropriate error.
		"scim_request_returns_400": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},

				PageSize: 1,
				Cursor:   "400",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 400.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		// This test case uses a random URL instead of the test server's URL, so we should expect an invalid cert.
		// This test case also verifies that the "https://" prefix is added to the address if it's not present
		// which is evident by the error message.
		"failed_to_make_get_page_request_invalid_certs": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: "example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},

				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute SCIM request: ` +
						`Get "https://example.com/Users?startIndex=1&count=1": tls: failed to verify certificate: x509: ` +
						`certificate signed by unknown authority.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"failed_to_make_get_page_request_invalid_host": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				// Deliberately add the extra "/".
				// This test case also indirectly verifies that the "https://" prefix is NOT added
				// if the prefix already exists.
				Address: "https:///example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},

				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute SCIM request: ` +
						`Get "https:///example.com/Users?startIndex=1&count=1": http: no Host in request URL.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		// The endpoint is setup to sleep for 5 seconds.
		// The request timeout is configured for 1 seconds.
		"user_request_timeout_with_common_config": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: baseURL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimUser,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "userName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "emails",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "type",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "primary",
									Type:       framework.AttributeTypeBool,
								},
							},
						},
					},
				},
				Config: &scim.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(1),
					},
				},
				PageSize: 1,
				Cursor:   "408",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to execute SCIM request: Get \"" + baseURL + "/Users?startIndex=408&count=1\": context deadline exceeded. Request exceeded configured timeout of 1 seconds. Please increase the request timeout.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},

		// The endpoint is setup to sleep for 5 seconds.
		// The request timeout is configured for 1 seconds.
		"group_request_timeout_with_common_config": {
			ctx: context.Background(),
			request: &framework.Request[scim.Config]{
				Address: baseURL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: testUsername,
						Password: testPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: scimGroup,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.meta.created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "members",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "type",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},
				Config: &scim.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(1),
					},
				},
				PageSize: 1,
				Cursor:   "408",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to execute SCIM request: Get \"" + baseURL + "/Groups?startIndex=408&count=1\": context deadline exceeded. Request exceeded configured timeout of 1 seconds. Please increase the request timeout.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
