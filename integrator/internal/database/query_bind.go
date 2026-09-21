package database

import (
	"errors"
	"strings"
	"unicode"
)

func (q *Query) bind(values ...any) *Query {

	if len(values)%2 != 0 {
		q.err = errors.New("bind() harus pasangan key,value")
		return q
	}

	params := make(map[string]any)

	for i := 0; i < len(values); i += 2 {

		key, ok := values[i].(string)

		if !ok {
			q.err = errors.New("parameter harus string")
			return q
		}

		params[key] = values[i+1]

	}

	var sql strings.Builder
	args := make([]any, 0)
	text := q.sql

	for i := 0; i < len(text); {
		if text[i] != ':' {
			sql.WriteByte(text[i])
			i++
			continue
		}

		i++

		start := i
		for i < len(text) {
			r := rune(text[i])

			if unicode.IsLetter(r) ||
				unicode.IsDigit(r) ||
				r == '_' {
				i++
				continue
			}

			break
		}

		name := text[start:i]

		value, ok := params[name]

		if !ok {
			q.err = errors.New("parameter :" + name + " belum di bind()")
			return q
		}

		sql.WriteByte('?')

		args = append(args, value)

	}

	q.sql = sql.String()
	q.args = args
	return q

}
