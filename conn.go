package pgsqlxx

import (
	"errors"
	"reflect"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx/reflectx"
)

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/sqlx.go#L26
var NameMapper = strings.ToLower
var origMapper = reflect.ValueOf(NameMapper)

// Rather than creating on init, this is created when necessary so that
// importers have time to customize the NameMapper.
var mpr *reflectx.Mapper

// mapper returns a valid mapper using the configured NameMapper func.
func mapper() *reflectx.Mapper {
	if mpr == nil {
		mpr = reflectx.NewMapperFunc("db", NameMapper)
	} else if origMapper != reflect.ValueOf(NameMapper) {
		// if NameMapper has changed, create a new mapper
		mpr = reflectx.NewMapperFunc("db", NameMapper)
		origMapper = reflect.ValueOf(NameMapper)
	}
	return mpr
}

func ConnectFromPool(pool *pgx.ConnPool) (*Connxx, error) {
	s := pool.Stat()

	if s.AvailableConnections <= 0 && s.CurrentConnections <= 0 {
		return nil, errors.New("no connections active in pool")
	}

	c, err := pool.Acquire()
	if err != nil {
		return nil, err
	}

	if !c.IsAlive() {
		pool.Release(c)
		deathErr := c.CauseOfDeath()
		if deathErr != nil {
			return nil, deathErr
		} else {
			return nil, errors.New("acquired dead connection from pool")
		}
	}

	pool.Release(c)

	return &Connxx{ConnPool: pool, unsafe: false, Mapper: mapper()}, nil

}

func MustConnectFromPool(pool *pgx.ConnPool) *Connxx {
	c, e := ConnectFromPool(pool)

	if e != nil {
		panic(e)
	}
	return c
}

type Connxx struct {
	*pgx.ConnPool
	unsafe bool
	Mapper *reflectx.Mapper
}

func (c *Connxx) Unsafe() *Connxx {
	c.unsafe = true
	return c
}

func (c *Connxx) MapperFunc(mf func(string) string) {
	c.Mapper = reflectx.NewMapperFunc("db", mf)
}

func (c *Connxx) Queryx(query string, args ...interface{}) (*Rows, error) {
	rows, err := c.ConnPool.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return &Rows{Rows: rows, unsafe: c.unsafe, Mapper: c.Mapper}, nil
}

func (c *Connxx) QueryRowx(query string, args ...interface{}) *Row {
	rows, err := c.Queryx(query, args...)
	return &Row{Rows: rows, err: err, unsafe: c.unsafe, Mapper: c.Mapper}
}

func (c *Connxx) Rebind(query string) string {
	return Rebind(query)
}

func (c *Connxx) Beginx() (*Tx, error) {
	tx, err := c.Begin()

	if err != nil {
		return nil, err
	}

	return &Tx{Tx: tx, unsafe: c.unsafe, Mapper: c.Mapper}, nil
}
