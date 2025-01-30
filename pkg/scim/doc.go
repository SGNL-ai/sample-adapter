// Copyright 2025 SGNL.ai, Inc.

/*
Package scim implements an adapter for the System for Cross-domain Identity Management (SCIM) protocol v2.

## Group membership

SCIM does not provide a dedicated endpoint for group membership data. Instead, it provides
- a `members` attribute on the Group resource i.e. a list of objects containing the user IDs the group contains
- a `groups` attribute on the User resource i.e. a list of  objects containing the group IDs the user is a member of

A typical SCIM server is likely to contain a relatively small number of groups compared to the number of users.
This means that the `members` attribute on the Group resource is likely to be very large,
and the `groups` attribute on the User resource is likely to be small.

As there's no pagination support for `members`/`groups`, the design decision is to
ingest group membership data from the User resource, which is relatively small
and to ignore the `members` attribute on the Group resource.

Group members are ingested as child entities on the SGNL Console.
*/
package scim
