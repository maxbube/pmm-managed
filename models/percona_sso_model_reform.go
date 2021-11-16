// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package models

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type perconaSSODetailsViewType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *perconaSSODetailsViewType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("percona_sso_details").
func (v *perconaSSODetailsViewType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *perconaSSODetailsViewType) Columns() []string {
	return []string{
		"client_id",
		"client_secret",
		"issuer_url",
		"scope",
		"created_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *perconaSSODetailsViewType) NewStruct() reform.Struct {
	return new(PerconaSSODetails)
}

// PerconaSSODetailsView represents percona_sso_details view or table in SQL database.
var PerconaSSODetailsView = &perconaSSODetailsViewType{
	s: parse.StructInfo{
		Type:    "PerconaSSODetails",
		SQLName: "percona_sso_details",
		Fields: []parse.FieldInfo{
			{Name: "ClientID", Type: "string", Column: "client_id"},
			{Name: "ClientSecret", Type: "string", Column: "client_secret"},
			{Name: "IssuerURL", Type: "string", Column: "issuer_url"},
			{Name: "Scope", Type: "string", Column: "scope"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
		},
		PKFieldIndex: -1,
	},
	z: new(PerconaSSODetails).Values(),
}

// String returns a string representation of this struct or record.
func (s PerconaSSODetails) String() string {
	res := make([]string, 5)
	res[0] = "ClientID: " + reform.Inspect(s.ClientID, true)
	res[1] = "ClientSecret: " + reform.Inspect(s.ClientSecret, true)
	res[2] = "IssuerURL: " + reform.Inspect(s.IssuerURL, true)
	res[3] = "Scope: " + reform.Inspect(s.Scope, true)
	res[4] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *PerconaSSODetails) Values() []interface{} {
	return []interface{}{
		s.ClientID,
		s.ClientSecret,
		s.IssuerURL,
		s.Scope,
		s.CreatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *PerconaSSODetails) Pointers() []interface{} {
	return []interface{}{
		&s.ClientID,
		&s.ClientSecret,
		&s.IssuerURL,
		&s.Scope,
		&s.CreatedAt,
	}
}

// View returns View object for that struct.
func (s *PerconaSSODetails) View() reform.View {
	return PerconaSSODetailsView
}

// check interfaces
var (
	_ reform.View   = PerconaSSODetailsView
	_ reform.Struct = (*PerconaSSODetails)(nil)
	_ fmt.Stringer  = (*PerconaSSODetails)(nil)
)

func init() {
	parse.AssertUpToDate(&PerconaSSODetailsView.s, new(PerconaSSODetails))
}