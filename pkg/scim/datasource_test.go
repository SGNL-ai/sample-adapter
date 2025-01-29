// Copyright 2025 SGNL.ai, Inc.
package scim_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/sample-adapter/pkg/scim"
	"github.com/sgnl-ai/sample-adapter/pkg/testutil"
)

var (
	scimUser  = "Users"
	scimGroup = "Groups"
)

func TestUserGetPage(t *testing.T) {
	scimClient := scim.NewClient(&http.Client{
		Timeout: time.Duration(60) * time.Second,
	})

	server := httptest.NewServer(TestServerHandler)
	defer server.Close()

	tests := map[string]struct {
		context context.Context
		request *scim.Request
		wantRes *scim.AdapterResponse
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &scim.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,

				EntityExternalID: scimUser,
				PageSize:         2,
			},
			wantRes: &scim.AdapterResponse{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "2819c223-7f76-453a-919d-413861904646", "userName": "Alex"},
					{"id": "c75ad752-64ae-4823-840d-ffa80929976c", "userName": "Bacong"},
				},
				NextCursor: "3",
			},
			wantErr: nil,
		},
		"middle_page": {
			context: context.Background(),
			request: &scim.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,

				EntityExternalID: scimUser,
				PageSize:         2,
				Cursor:           "3",
			},
			wantRes: &scim.AdapterResponse{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "e2be737c-61f5-4abe-8797-1e816b15cec8", "userName": "Carol"},
					{"id": "89fa657e-3ef5-49e3-bb34-b3255e04a8bb", "userName": "David"},
				},
				NextCursor: "5",
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &scim.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,

				EntityExternalID: scimUser,
				PageSize:         2,
				Cursor:           "5",
			},
			wantRes: &scim.AdapterResponse{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"schemas": []interface{}{
							"urn:ietf:params:scim:schemas:core:2.0:User",
							"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
						},
						"id":         "2819c223-7f76-453a-919d-413861904000",
						"externalId": "701984",
						"userName":   "bjensen@example.com",
						"name": map[string]interface{}{
							"formatted":       "Ms. Barbara J Jensen, III",
							"familyName":      "Jensen",
							"givenName":       "Barbara",
							"middleName":      "Jane",
							"honorificPrefix": "Ms.",
							"honorificSuffix": "III",
						},
						"displayName": "Babs Jensen",
						"nickName":    "Babs",
						"profileUrl":  "https://login.example.com/bjensen",
						"emails": []interface{}{
							map[string]interface{}{
								"value":   "bjensen@example.com",
								"display": "bjensen@example.com",
								"type":    "work",
								"primary": true,
							},
							map[string]interface{}{
								"value": "babs@jensen.org",
								"type":  "home",
							},
						},
						"addresses": []interface{}{
							map[string]interface{}{
								"streetAddress": "100 Universal City Plaza",
								"locality":      "Hollywood",
								"region":        "CA",
								"postalCode":    "91608",
								"country":       "USA",
								"formatted":     "100 Universal City Plaza\nHollywood, CA 91608 USA",
								"type":          "work",
							},
							map[string]interface{}{
								"streetAddress": "456 Hollywood Blvd",
								"locality":      "Hollywood",
								"region":        "CA",
								"postalCode":    "91608",
								"country":       "USA",
								"formatted":     "456 Hollywood Blvd\nHollywood, CA 91608 USA",
								"type":          "home",
							},
						},
						"phoneNumbers": []interface{}{
							map[string]interface{}{
								"value":   "555-555-5555",
								"display": "555-555-5555",
								"primary": true,
								"type":    "work",
							},
							map[string]interface{}{
								"value": "555-555-4444",
								"type":  "mobile",
							},
						},
						"ims": []interface{}{
							map[string]interface{}{
								"value":   "someaimhandle",
								"display": "someaimhandle",
								"type":    "aim",
								"primary": true,
							},
						},
						"photos": []interface{}{
							map[string]interface{}{
								"value":   "https://photos.example.com/profilephoto/72930000000Ccne/F",
								"display": "https://photos.example.com/profilephoto/72930000000Ccne/F",
								"type":    "photo",
								"primary": true,
							},
							map[string]interface{}{
								"value": "https://photos.example.com/profilephoto/72930000000Ccne/T",
								"type":  "thumbnail",
							},
						},
						"userType":          "Employee",
						"title":             "Tour Guide",
						"preferredLanguage": "en-US",
						"locale":            "en-US",
						"timezone":          "America/Los_Angeles",
						"active":            true,
						"password":          "t1meMa$heen",
						"groups": []interface{}{
							map[string]interface{}{
								"value":   "e9e30dba-f08f-4109-8486-d5c6a331660a",
								"$ref":    "../Groups/e9e30dba-f08f-4109-8486-d5c6a331660a",
								"display": "Tour Guides",
								"type":    "direct",
							},
							map[string]interface{}{
								"value":   "fc348aa8-3835-40eb-a20b-c726e15c55b5",
								"$ref":    "../Groups/fc348aa8-3835-40eb-a20b-c726e15c55b5",
								"display": "Employees",
							},
							map[string]interface{}{
								"value":   "71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
								"$ref":    "../Groups/71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
								"display": "US Employees",
							},
						},
						"entitlements": []interface{}{
							map[string]interface{}{
								"value":   "e9e30dba-f08f-4109-8486-d5c6a331abc",
								"display": "E1",
								"type":    "entitlement",
								"primary": true,
							},
						},
						"roles": []interface{}{
							map[string]interface{}{
								"value":   "e9e30dba-f08f-4109-8486-d5c6a33role",
								"display": "Role A",
								"type":    "role",
								"primary": true,
							},
						},
						"x509Certificates": []interface{}{
							map[string]interface{}{
								"value":   "some_certificate_value",
								"display": "SGNLCertificate",
								"type":    "secret",
								"primary": true,
							},
						},
						"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
							"employeeNumber": "701984",
							"costCenter":     "4130",
							"organization":   "Universal Studios",
							"division":       "Theme Park",
							"department":     "Tour Operations",
							"manager": map[string]interface{}{
								"value":       "26118915-6090-4610-87e4-49d8ca9f808d",
								"$ref":        "../Users/26118915-6090-4610-87e4-49d8ca9f808d",
								"displayName": "John Smith",
							},
						},
						"meta": map[string]interface{}{
							"resourceType": "User",
							"created":      "2010-01-23T04:56:22Z",
							"lastModified": "2011-05-13T04:42:34Z",
							"version":      "W/\"3694e05e9dff591\"",
							"location":     "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := scimClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %+v, wantRes: %+v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGroupGetPage(t *testing.T) {
	scimClient := scim.NewClient(&http.Client{
		Timeout: time.Duration(60) * time.Second,
	})

	server := httptest.NewServer(TestServerHandler)
	defer server.Close()

	tests := map[string]struct {
		context context.Context
		request *scim.Request
		wantRes *scim.AdapterResponse
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &scim.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,

				EntityExternalID: scimGroup,
				PageSize:         2,
			},
			wantRes: &scim.AdapterResponse{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":          "c3a26dd3-27a0-4dec-a2ac-ce211e105f97",
						"displayName": "Group A",
						"members": []interface{}{
							map[string]interface{}{
								"value":   "6c5bb468-14b2-4183-baf2-06d523e03bd3",
								"$ref":    "https://example.com/v2/Groups/6c5bb468-14b2-4183-baf2-06d523e03bd3",
								"display": "Group",
							},
						},
					},
					{
						"id":          "6c5bb468-14b2-4183-baf2-06d523e03bd3",
						"displayName": "Group B",
					},
				},
				NextCursor: "3",
			},
			wantErr: nil,
		},
		"middle_page": {
			context: context.Background(),
			request: &scim.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,

				EntityExternalID: scimGroup,
				PageSize:         2,
				Cursor:           "3",
			},
			wantRes: &scim.AdapterResponse{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":          "e2be737c-61f5-4abe-8797-1e816b15cec8",
						"displayName": "Group C",
						"members": []interface{}{
							map[string]interface{}{
								"value":   "89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
								"$ref":    "https://example.com/v2/Groups/89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
								"display": "Group",
							},
						},
					},
					{
						"id":          "89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
						"displayName": "Group D",
					},
				},
				NextCursor: "5",
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &scim.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,

				EntityExternalID: scimGroup,
				PageSize:         2,
				Cursor:           "5",
			},
			wantRes: &scim.AdapterResponse{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":          "e9e30dba-f08f-4109-8486-d5c6a331660a",
						"displayName": "Tour Guides",
						"schemas": []interface{}{
							"urn:ietf:params:scim:schemas:core:2.0:Group",
						},
						"members": []interface{}{
							map[string]interface{}{
								"value": "2819c223-7f76-453a-919d-413861904646",
								"$ref":  "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
								"type":  "User",
							},
							map[string]interface{}{
								"value": "6c5bb468-14b2-4183-baf2-06d523e03bd3",
								"$ref":  "https://example.com/v2/Groups/6c5bb468-14b2-4183-baf2-06d523e03bd3",
								"type":  "Group",
							},
						},
						"meta": map[string]interface{}{
							"resourceType": "Group",
							"created":      "2010-01-23T04:56:22Z",
							"lastModified": "2011-05-13T04:42:34Z",
							"version":      "W/\"3694e05e9dff592\"",
							"location":     "https://example.com/v2/Groups/e9e30dba-f08f-4109-8486-d5c6a331660a",
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := scimClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %+v, wantRes: %+v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestConstructURL(t *testing.T) {
	tests := map[string]struct {
		request *scim.Request
		cursor  string
		wantURL string
	}{
		"users": {
			request: &scim.Request{
				BaseURL: "https://scim.com",

				PageSize:         10,
				EntityExternalID: scimUser,
			},
			cursor:  "10",
			wantURL: "https://scim.com/Users?startIndex=10&count=10",
		},
		"groups": {
			request: &scim.Request{
				BaseURL: "https://scim.com",

				PageSize:         10,
				EntityExternalID: scimGroup,
			},
			cursor:  "10",
			wantURL: "https://scim.com/Groups?startIndex=10&count=10",
		},
		"userQueryParams": {
			request: &scim.Request{
				BaseURL: "https://scim.com",

				PageSize:         10,
				EntityExternalID: scimUser,
				QueryParams: scim.QueryParams{
					SortBy:    "userName",
					Ascending: testutil.GenPtr(true),
					Filter:    `userType eq "Employee" and (emails co "example.com" or emails.value co "example.org")`,
				},
			},
			cursor: "10",
			wantURL: `https://scim.com/Users?` +
				`startIndex=10&` +
				`count=10&` +
				`filter=userType+eq+%22Employee%22+and+%28emails+co+%22example.com%22+or+emails.value+co+%22example.org%22%29&` +
				`sortBy=userName&` +
				`sortOrder=ascending`,
		},
		"groupQueryParams": {
			request: &scim.Request{
				BaseURL: "https://scim.com",

				PageSize:         10,
				EntityExternalID: scimGroup,
				QueryParams: scim.QueryParams{
					SortBy:    "displayName",
					Ascending: testutil.GenPtr(true),
					Filter:    `displayName eq "SGNL"`,
				},
			},
			cursor: "10",
			wantURL: `https://scim.com/Groups?` +
				`startIndex=10&` +
				`count=10&` +
				`filter=displayName+eq+%22SGNL%22&` +
				`sortBy=displayName&` +
				`sortOrder=ascending`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotURL := scim.GenerateURL(
				tt.request.BaseURL,
				tt.request.EntityExternalID,
				tt.request.PageSize,
				tt.cursor,
				tt.request.QueryParams,
			)

			if !reflect.DeepEqual(gotURL, tt.wantURL) {
				t.Errorf("gotURL: %v, wantURL: %v", gotURL, tt.wantURL)
			}
		})
	}
}
