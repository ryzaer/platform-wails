package database

func (q *Query) Exec(values ...any) error {

	if len(values) > 0 {
		q.bind(values...)
	}

	if q.err != nil {
		return q.err
	}

	_, err := q.db.conn.Exec(
		q.sql,
		q.args...,
	)

	return err
}

func (q *Query) LastInsertID(values ...any) (int64, error) {

	if len(values) > 0 {
		q.bind(values...)
	}

	if q.err != nil {
		return 0, q.err
	}

	result, err := q.db.conn.Exec(
		q.sql,
		q.args...,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (q *Query) RowsAffected(values ...any) (int64, error) {

	if len(values) > 0 {
		q.bind(values...)
	}

	if q.err != nil {
		return 0, q.err
	}

	result, err := q.db.conn.Exec(
		q.sql,
		q.args...,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
