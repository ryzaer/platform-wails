package database

import "database/sql"

type Query struct {
	db   *Connect
	sql  string
	args []any
	err  error
}

func (q *Query) Error() error {
	return q.err
}

func (q *Query) Row() *sql.Row {

	if q.err != nil {
		return nil
	}

	return q.db.conn.QueryRow(
		q.sql,
		q.args...,
	)

}

func (q *Query) Rows() (*sql.Rows, error) {

	if q.err != nil {
		return nil, q.err
	}

	return q.db.conn.Query(
		q.sql,
		q.args...,
	)

}
