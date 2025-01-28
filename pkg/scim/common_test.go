// Copyright 2025 SGNL.ai, Inc.
package scim_test

import (
	"net/http"
	"time"
)

// Define the endpoints and responses for the mock SCIM server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.RequestURI() {

	// User endpoints
	case "/Users?startIndex=1&count=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				"urn:ietf:params:scim:api:messages:2.0:ListResponse"
			],
			"totalResults": 5,
			"itemsPerPage": 2,
			"startIndex": 1,
			"Resources": [
				{
					"id": "2819c223-7f76-453a-919d-413861904646",
					"userName": "Alex"
				},
				{
					"id": "c75ad752-64ae-4823-840d-ffa80929976c","userName": "Bacong"
				}
			]
		}`))
	case "/Users?startIndex=3&count=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas":["urn:ietf:params:scim:api:messages:2.0:ListResponse"],
			"totalResults":5,
			"itemsPerPage":2,
			"startIndex":3,
     		"Resources":[
			  {
				"id":"e2be737c-61f5-4abe-8797-1e816b15cec8",
				"userName":"Carol"
			  },
			  {
				"id":"89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
				"userName":"David"
			  }
			]
		  }`))
	case "/Users?startIndex=5&count=2":
		// Response is a Full Enterprise User Extension Representation
		// https://datatracker.ietf.org/doc/html/rfc7643#section-8.3
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				"urn:ietf:params:scim:api:messages:2.0:ListResponse"
			],
			"totalResults": 5,
			"itemsPerPage": 2,
			"startIndex": 5,
			"Resources": [
				{
					"schemas": [
						"urn:ietf:params:scim:schemas:core:2.0:User",
						"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
					],
					"id": "2819c223-7f76-453a-919d-413861904000",
					"externalId": "701984",
					"userName": "bjensen@example.com",
					"name": {
						"formatted": "Ms. Barbara J Jensen, III",
						"familyName": "Jensen",
						"givenName": "Barbara",
						"middleName": "Jane",
						"honorificPrefix": "Ms.",
						"honorificSuffix": "III"
					},
					"displayName": "Babs Jensen",
					"nickName": "Babs",
					"profileUrl": "https://login.example.com/bjensen",
					"emails": [
						{
							"value": "bjensen@example.com",
							"display": "bjensen@example.com",
							"type": "work",
							"primary": true
						},
						{
							"value": "babs@jensen.org",
							"type": "home"
						}
					],
					"addresses": [
						{
							"streetAddress": "100 Universal City Plaza",
							"locality": "Hollywood",
							"region": "CA",
							"postalCode": "91608",
							"country": "USA",
							"formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
							"type": "work"
						},
						{
							"streetAddress": "456 Hollywood Blvd",
							"locality": "Hollywood",
							"region": "CA",
							"postalCode": "91608",
							"country": "USA",
							"formatted": "456 Hollywood Blvd\nHollywood, CA 91608 USA",
							"type": "home"
						}
					],
					"phoneNumbers": [
						{
							"value": "555-555-5555",
							"display": "555-555-5555",
							"type": "work",
							"primary": true
						},
						{
							"value": "555-555-4444",
							"type": "mobile"
						}
					],
					"ims": [
						{
							"value": "someaimhandle",
							"display": "someaimhandle",
							"type": "aim",
							"primary": true
						}
					],
					"photos": [
						{
							"value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
							"display": "https://photos.example.com/profilephoto/72930000000Ccne/F",
							"type": "photo",
							"primary": true
						},
						{
							"value": "https://photos.example.com/profilephoto/72930000000Ccne/T",
							"type": "thumbnail"
						}
					],
					"userType": "Employee",
					"title": "Tour Guide",
					"preferredLanguage": "en-US",
					"locale": "en-US",
					"timezone": "America/Los_Angeles",
					"active": true,
					"password": "t1meMa$heen",
					"groups": [
						{
							"value": "e9e30dba-f08f-4109-8486-d5c6a331660a",
							"$ref": "../Groups/e9e30dba-f08f-4109-8486-d5c6a331660a",
							"display": "Tour Guides",
							"type": "direct"
						},
						{
							"value": "fc348aa8-3835-40eb-a20b-c726e15c55b5",
							"$ref": "../Groups/fc348aa8-3835-40eb-a20b-c726e15c55b5",
							"display": "Employees"
						},
						{
							"value": "71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
							"$ref": "../Groups/71ddacd2-a8e7-49b8-a5db-ae50d0a5bfd7",
							"display": "US Employees"
						}
					],
					"entitlements": [
						{
							"value": "e9e30dba-f08f-4109-8486-d5c6a331abc",
							"display": "E1",
							"type": "entitlement",
							"primary": true
						}
					],
					"roles": [
						{
							"value": "e9e30dba-f08f-4109-8486-d5c6a33role",
							"display": "Role A",
							"type": "role",
							"primary": true
						}
					],
					"x509Certificates": [
						{
							"value": "some_certificate_value",
							"display": "SGNLCertificate",
							"type": "secret",
							"primary": true
						}
					],
					"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": {
						"employeeNumber": "701984",
						"costCenter": "4130",
						"organization": "Universal Studios",
						"division": "Theme Park",
						"department": "Tour Operations",
						"manager": {
							"value": "26118915-6090-4610-87e4-49d8ca9f808d",
							"$ref": "../Users/26118915-6090-4610-87e4-49d8ca9f808d",
							"displayName": "John Smith"
						}
					},
					"meta": {
						"resourceType": "User",
						"created": "2010-01-23T04:56:22Z",
						"lastModified": "2011-05-13T04:42:34Z",
						"version": "W\/\"3694e05e9dff591\"",
						"location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
					}
				}
			]
		}`))

	// Group endpoints
	// Group B is a member of Group A. Group D is a member of Group C.
	// Members of group "Tour Guides" is a user and a group.
	case "/Groups?startIndex=1&count=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				"urn:ietf:params:scim:api:messages:2.0:ListResponse"
			],
			"totalResults": 5,
			"itemsPerPage": 2,
			"startIndex": 1,
			"Resources": [
				{
					"id": "c3a26dd3-27a0-4dec-a2ac-ce211e105f97",
					"displayName": "Group A",
					"members": [
						{
							"value": "6c5bb468-14b2-4183-baf2-06d523e03bd3",
							"$ref": "https://example.com/v2/Groups/6c5bb468-14b2-4183-baf2-06d523e03bd3",
							"display": "Group"
						}
					]
				},
				{
					"id": "6c5bb468-14b2-4183-baf2-06d523e03bd3",
					"displayName": "Group B"
				}
			]
		}`))
	case "/Groups?startIndex=3&count=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				"urn:ietf:params:scim:api:messages:2.0:ListResponse"
			],
			"totalResults": 5,
			"itemsPerPage": 2,
			"startIndex": 3,
			"Resources": [
				{
					"id": "e2be737c-61f5-4abe-8797-1e816b15cec8",
					"displayName": "Group C",
					"members": [
						{
							"value": "89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
							"$ref": "https://example.com/v2/Groups/89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
							"display": "Group"
						}
					]
				},
				{
					"id": "89fa657e-3ef5-49e3-bb34-b3255e04a8bb",
					"displayName": "Group D"
				}
			]
		}`))
	case "/Groups?startIndex=5&count=2":
		// https://datatracker.ietf.org/doc/html/rfc7643#section-8.4
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				"urn:ietf:params:scim:api:messages:2.0:ListResponse"
			],
			"totalResults": 5,
			"itemsPerPage": 2,
			"startIndex": 5,
			"Resources": [
				{
					"schemas": [
						"urn:ietf:params:scim:schemas:core:2.0:Group"
					],
					"id": "e9e30dba-f08f-4109-8486-d5c6a331660a",
					"displayName": "Tour Guides",
					"members": [
						{
							"value": "2819c223-7f76-453a-919d-413861904646",
							"$ref": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646",
							"type": "User"
						},
						{
							"value": "6c5bb468-14b2-4183-baf2-06d523e03bd3",
							"$ref": "https://example.com/v2/Groups/6c5bb468-14b2-4183-baf2-06d523e03bd3",
							"type": "Group"
						}
					],
					"meta": {
						"resourceType": "Group",
						"created": "2010-01-23T04:56:22Z",
						"lastModified": "2011-05-13T04:42:34Z",
						"version": "W\/\"3694e05e9dff592\"",
						"location": "https://example.com/v2/Groups/e9e30dba-f08f-4109-8486-d5c6a331660a"
					}
				}
			]
		}`))

	// Additional endpoints to facilitate testing
	// Simulate a bad request
	case "/Users?startIndex=400&count=1":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"err": "bad request"
		}`))

	// Simulate a user request timeout
	case "/Users?startIndex=408&count=1":
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"err": "This is meant to simulate a long running task resulting in a timeout."
		}`))

	// Simulate a group request timeout
	case "/Groups?startIndex=408&count=1":
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"err": "This is meant to simulate a long running task resulting in a timeout."
		}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
