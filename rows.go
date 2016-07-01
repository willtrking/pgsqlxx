package pgsqlxx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

// Allow easy access to pgsqlxx Rows without needing to use the rest of the library
func RowsFromRows(rows *pgx.Rows, mapper *reflectx.Mapper) *Rows {
	return &Rows{
		Rows:   rows,
		Mapper: mapper,
	}
}

// Allow easy access to pgsqlxx Rows without needing to use the rest of the library
func RowsFromRowsUnsafe(rows *pgx.Rows, mapper *reflectx.Mapper) *Rows {
	return &Rows{
		Rows:   rows,
		unsafe: true,
		Mapper: mapper,
	}
}

type Rows struct {
	*pgx.Rows
	Mapper *reflectx.Mapper

	unsafe  bool
	started bool
	fields  [][]int
	values  []interface{}
}

func (r *Rows) MapScan(dest map[string]interface{}) error {
	return sqlx.MapScan(r, dest)
}

func (r *Rows) SliceScan() ([]interface{}, error) {
	return sqlx.SliceScan(r)
}

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L560
func (r *Rows) StructScan(dest interface{}) error {
	rowErr := r.Err()
	if rowErr != nil {
		return rowErr
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

func (r *Rows) Columns() ([]string, error) {
	fieldDescriptions := r.Rows.FieldDescriptions()
	names := make([]string, 0, len(fieldDescriptions))
	for _, fd := range fieldDescriptions {
		names = append(names, fd.Name)
	}
	return names, nil
}

func (r *Rows) Closex() error {
	r.Close()
	return r.Err()
}
