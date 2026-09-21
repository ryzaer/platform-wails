package material

import (
	"fmt"

	"app-platform/internal/database"
)

type Option struct {
	Value any    `json:"value"`
	Label string `json:"label"`
}

func SupplierOptions(db *database.Connect) ([]Option, error) {
	rows, err := db.Query(`
		SELECT id, name
		FROM suppliers
		ORDER BY name
	`).FetchAll()

	if err != nil {
		return nil, err
	}

	out := make([]Option, 0, len(rows))

	for _, row := range rows {
		out = append(out, Option{
			Value: row["id"],
			Label: toString(row["name"]),
		})
	}

	return out, nil
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	if b, ok := v.([]byte); ok {
		return string(b)
	}

	return fmt.Sprint(v)
}
