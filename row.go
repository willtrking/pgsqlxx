package pgsqlxx

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx/reflectx"
)

// Allow easy access to pgsqlxx Row without needing to use the rest of the library
func RowFromRows(rows *pgx.Rows, mapper *reflectx.Mapper) *Row {
	return &Row{
		Rows: &Rows{
			Rows:   rows,
			Mapper: mapper,
		},
		Mapper: mapper,
		err:    nil,
	}
}

// Allow easy access to pgsqlxx Row without needing to use the rest of the library
func RowFromRowsUnsafe(rows *pgx.Rows, mapper *reflectx.Mapper) *Row {
	return &Row{
		Rows: &Rows{
			Rows:   rows,
			Mapper: mapper,
		},
		Mapper: mapper,
		err:    nil,
		unsafe: true,
	}
}

// Allow easy access to pgsqlxx Row without needing to use the rest of the library
func RowFromRowsx(rows *Rows, mapper *reflectx.Mapper) *Row {
	return &Row{
		Rows:   rows,
		Mapper: mapper,
		err:    nil,
	}
}

// Allow easy access to pgsqlxx Row without needing to use the rest of the library
func RowFromRowsUnsafex(rows *Rows, mapper *reflectx.Mapper) *Row {
	return &Row{
		Rows:   rows,
		Mapper: mapper,
		err:    nil,
		unsafe: true,
	}
}

type Row struct {
	*Rows
	err    error
	unsafe bool
	Mapper *reflectx.Mapper
}

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L173
func (r *Row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	// TODO(bradfitz): for now we need to defensively clone all
	// []byte that the driver returned (not permitting
	// *RawBytes in Rows.Scan), since we're about to close
	// the Rows in our defer, when we return from this function.
	// the contract with the driver.Next(...) interface is that it
	// can return slices into read-only temporary memory that's
	// only valid until the next Scan/Close.  But the TODO is that
	// for a lot of drivers, this copy will be unnecessary.  We
	// should provide an optional interface for drivers to
	// implement to say, "don't worry, the []bytes that I return
	// from Next will not be modified again." (for instance, if
	// they were obtained from the network anyway) But for now we
	// don't care.
	defer r.Close()
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.Next() {
		if err := r.Err(); err != nil {
			return err
		}
		return pgx.ErrNoRows
	}
	err := r.Rows.Scan(dest...)
	if err != nil {
		return err
	}
	// Make sure the query can be processed to completion with no errors.
	r.Close()

	if err := r.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Row) StructScan(dest interface{}) error {

	if r.err != nil {
		return r.err
	}

	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = reflect.Indirect(v)

	if !r.started {
		columns, err := r.Columns()
		if err != nil {
			return err
		}
		m := r.Mapper

		r.fields = m.TraversalsByName(v.Type(), columns)

		// if we are not unsafe and are missing fields, return an error
		if !r.unsafe {
			if f, err := missingFields(r.fields); err != nil {
				return fmt.Errorf("missing destination name %s", columns[f])
			}
		}

		r.values = make([]interface{}, len(columns))
		r.started = true
	}

	err := fieldsByTraversal(v, r.fields, r.values, true)
	if err != nil {
		return err
	}

	// scan into the struct field pointers and append to our results
	err = r.Scan(r.values...)
	if err != nil {
		return err
	}

	return r.Err()

}

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L217
func (r *Row) Columns() ([]string, error) {
	if r.err != nil {
		return []string{}, r.err
	}
	return r.Rows.Columns()
}

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L225
func (r *Row) Err() error {
	return r.err
}
