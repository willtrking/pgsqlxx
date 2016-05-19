package pgsqlxx

import (
	"errors"

	"github.com/jackc/pgx"
)

func ResultFromCommandTag(t pgx.CommandTag) *Result {
	return &Result{affected: t.RowsAffected(), err: nil}
}

func ResultFromExec(t pgx.CommandTag, err error) *Result {
	if err != nil {
		return &Result{affected: 0, err: err}
	}
	return &Result{affected: t.RowsAffected(), err: nil}

}

// Make it fit with database/sql Result interface
type Result struct {
	affected int64
	err      error
}

func (r *Result) LastInsertId() (int64, error) {
	return 0, errors.New("pgx does not support LastInsertId, try a Query with RETURNING")
}

func (r *Result) RowsAffected() (int64, error) {
	return r.affected, r.err
}

func (r *Result) RowsAffectedx() int64 {
	return r.affected
}

func (r *Result) LastInsertIdx() int64 {
	return 0
}

func (r *Result) RowsAffectedErr() error {
	return r.err
}

func (r *Result) LastInsertIdErr() error {
	return errors.New("pgx does not support LastInsertId, try a Query with RETURNING")
}
