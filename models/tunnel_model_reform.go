// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package models

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type tunnelTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *tunnelTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("tunnels").
func (v *tunnelTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *tunnelTableType) Columns() []string {
	return []string{
		"tunnel_id",
		"tunnel_type",
		"pmm_agent_id",
		"connect_port",
		"listen_port",
		"created_at",
		"updated_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *tunnelTableType) NewStruct() reform.Struct {
	return new(Tunnel)
}

// NewRecord makes a new record for that table.
func (v *tunnelTableType) NewRecord() reform.Record {
	return new(Tunnel)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *tunnelTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// TunnelTable represents tunnels view or table in SQL database.
var TunnelTable = &tunnelTableType{
	s: parse.StructInfo{
		Type:    "Tunnel",
		SQLName: "tunnels",
		Fields: []parse.FieldInfo{
			{Name: "TunnelID", Type: "string", Column: "tunnel_id"},
			{Name: "TunnelType", Type: "TunnelType", Column: "tunnel_type"},
			{Name: "PMMAgentID", Type: "string", Column: "pmm_agent_id"},
			{Name: "ConnectPort", Type: "uint16", Column: "connect_port"},
			{Name: "ListenPort", Type: "uint16", Column: "listen_port"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(Tunnel).Values(),
}

// String returns a string representation of this struct or record.
func (s Tunnel) String() string {
	res := make([]string, 7)
	res[0] = "TunnelID: " + reform.Inspect(s.TunnelID, true)
	res[1] = "TunnelType: " + reform.Inspect(s.TunnelType, true)
	res[2] = "PMMAgentID: " + reform.Inspect(s.PMMAgentID, true)
	res[3] = "ConnectPort: " + reform.Inspect(s.ConnectPort, true)
	res[4] = "ListenPort: " + reform.Inspect(s.ListenPort, true)
	res[5] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[6] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Tunnel) Values() []interface{} {
	return []interface{}{
		s.TunnelID,
		s.TunnelType,
		s.PMMAgentID,
		s.ConnectPort,
		s.ListenPort,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Tunnel) Pointers() []interface{} {
	return []interface{}{
		&s.TunnelID,
		&s.TunnelType,
		&s.PMMAgentID,
		&s.ConnectPort,
		&s.ListenPort,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *Tunnel) View() reform.View {
	return TunnelTable
}

// Table returns Table object for that record.
func (s *Tunnel) Table() reform.Table {
	return TunnelTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Tunnel) PKValue() interface{} {
	return s.TunnelID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Tunnel) PKPointer() interface{} {
	return &s.TunnelID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Tunnel) HasPK() bool {
	return s.TunnelID != TunnelTable.z[TunnelTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.TunnelID = pk.
func (s *Tunnel) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = TunnelTable
	_ reform.Struct = (*Tunnel)(nil)
	_ reform.Table  = TunnelTable
	_ reform.Record = (*Tunnel)(nil)
	_ fmt.Stringer  = (*Tunnel)(nil)
)

func init() {
	parse.AssertUpToDate(&TunnelTable.s, new(Tunnel))
}
