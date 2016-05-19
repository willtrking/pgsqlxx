package pgsqlxx

import (
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx/reflectx"
)

type Tx struct {
	*pgx.Tx
	unsafe bool
	Mapper *reflectx.Mapper
}

func (t *Tx) Queryx(query string, args ...interface{}) (*Rows, error) {
	rows, err := t.Tx.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return &Rows{Rows: rows, unsafe: t.unsafe, Mapper: t.Mapper}, nil
}

func (t *Tx) QueryRowx(query string, args ...interface{}) *Row {
	rows, err := t.Queryx(query, args...)
	return &Row{Rows: rows, err: err, unsafe: t.unsafe, Mapper: t.Mapper}
}

func (t *Tx) Rebind(query string) string {
	return Rebind(query)
}
