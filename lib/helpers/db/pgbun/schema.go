package pgbun

import (
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Sort struct {
	SortKey         string
	Ascending       bool
	CaseInsensitive bool
}

type BaseSchema struct {
	ID        uuid.UUID `bun:"id,pk,type:uuid,nullzero" pg:"id,pk,type:uuid" json:"_id"`
	CreatedAt time.Time `bun:"created_at,pk,type:time,nullzero" pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,pk,type:time,nullzero" pg:"updated_at" json:"updated_at"`
}

type TableName bun.Ident

func (t TableName) Ident() bun.Ident {
	return bun.Ident(t)
}

func (t TableName) All() ColumnName {
	return NewColumnName(t, "*")
}

type ColumnName struct {
	value string
	table TableName
}

func NewColumnName(tableName TableName, val string) ColumnName {
	return ColumnName{
		value: val,
		table: tableName,
	}
}

func (c ColumnName) Ident() bun.Ident {
	return bun.Ident(c.value)
}

func (c ColumnName) String() string {
	return c.value
}

func (c ColumnName) WithTable(customTable ...TableName) bun.Ident {
	table := c.table
	if len(customTable) > 0 {
		table = customTable[0]
	}

	return bun.Ident(string(table) + "." + c.value)
}
