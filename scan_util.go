package pgsqlxx

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L964
func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L985
func missingFields(transversals [][]int) (field int, err error) {
	for i, t := range transversals {
		if len(t) == 0 {
			return i, errors.New("missing field")
		}
	}
	return 0, nil
}
