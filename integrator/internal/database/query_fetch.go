package database

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("record not found")

func (q *Query) FetchAll(values ...any) ([]map[string]any, error) {

	if len(values) > 0 {
		q.bind(values...)
	}

	if q.err != nil {
		return nil, q.err
	}

	rows, err := q.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, 0)

	valuesRow := make([]any, len(columns))
	ptr := make([]any, len(columns))

	for i := range valuesRow {
		ptr[i] = &valuesRow[i]
	}

	for rows.Next() {

		if err := rows.Scan(ptr...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(columns))

		for i, col := range columns {

			switch v := valuesRow[i].(type) {
			case []byte:
				row[col] = string(v)
			default:
				row[col] = v
			}

		}

		result = append(result, row)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (q *Query) Fetch(values ...any) (map[string]any, error) {

	rows, err := q.FetchAll(values...)

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, ErrNotFound
	}

	return rows[0], nil

}

// ---------------------------------------------

func (q *Query) FetchGroup(groups ...string) (map[string]any, error) {

	rows, err := q.FetchAll()
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, errors.New("FetchAll requires at least one group field")
	}

	return group(rows, groups), nil

}

// ---------------------------------------------

func group(rows []map[string]any, groups []string) map[string]any {
	result := make(map[string]any)
	key := groups[0]

	// level terakhir
	if len(groups) == 1 {

		for _, row := range rows {

			val, ok := row[key]
			if !ok {
				continue
			}

			groupKey := fmt.Sprint(val)

			delete(row, key)

			list, _ := result[groupKey].([]map[string]any)
			list = append(list, row)

			result[groupKey] = list
		}

		return result

	}

	// masih ada level berikutnya

	bucket := make(map[string][]map[string]any)
	for _, row := range rows {

		val, ok := row[key]
		if !ok {
			continue
		}

		groupKey := fmt.Sprint(val)

		delete(row, key)

		bucket[groupKey] = append(bucket[groupKey], row)

	}

	for k, v := range bucket {
		result[k] = group(v, groups[1:])
	}

	return result

}
