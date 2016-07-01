package pgsqlxx

import (
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx/reflectx"
)

// Allow easy access to pgsqlxx Tx without needing to use the rest of the library
func TxFromTx(tx *pgx.Tx, mapper *reflectx.Mapper) *Tx {
	return &Tx{Tx: tx, unsafe: false, Mapper: mapper}
}

// Allow easy access to pgsqlxx Tx without needing to use the rest of the library
func TxFromTxUnsafe(tx *pgx.Tx, mapper *reflectx.Mapper) *Tx {
	return &Tx{Tx: tx, unsafe: true, Mapper: mapper}
}

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

	r := RowsFromRows(rows, t.Mapper)
	r.unsafe = t.unsafe

	return r, nil
}

func (t *Tx) QueryRowx(query string, args ...interface{}) *Row {
	rows, err := t.Queryx(query, args...)

	r := RowFromRowsx(rows, t.Mapper)
	r.unsafe = t.unsafe
	r.err = err

	return r
}

func (t *Tx) Rebind(query string) string {
	return Rebind(query)
}

func (t *Tx) Execx(sql string, args ...interface{}) (*Result, error) {
	r, err := t.Exec(sql, args...)
	if err != nil {
		return nil, err
	}
	return ResultFromCommandTag(r), nil
}

func (t *Tx) Unsafe() *Tx {
	t.unsafe = true
	return t
}

func (t *Tx) IsTx() bool {
	return true
}

func (t *Tx) DriverName() string {
	return "pgx"
}
