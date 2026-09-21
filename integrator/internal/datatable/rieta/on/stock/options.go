package stock

import (
	"app-platform/internal/database"
	"fmt"
)

type Option struct {
	Value any    `json:"value"`
	Label string `json:"label"`
}

func MaterialOptions(db *database.Connect) ([]Option, error) {
	rows, err := db.Query(`SELECT id, item_name FROM material WHERE status = 1 ORDER BY item_name`).FetchAll()
	if err != nil {
		return nil, err
	}
	out := make([]Option, 0, len(rows))
	for _, row := range rows {
		out = append(out, Option{Value: row["id"], Label: toString(row["item_name"])})
	}
	return out, nil
}

func SizeOptions(db *database.Connect) ([]Option, error) {
	rows, err := db.Query(`SELECT id, name FROM sizes ORDER BY name`).FetchAll()
	if err != nil {
		return nil, err
	}
	out := make([]Option, 0, len(rows))
	for _, row := range rows {
		out = append(out, Option{Value: row["id"], Label: toString(row["name"])})
	}
	return out, nil
}

func BranchOptions(db *database.Connect) ([]Option, error) {
	rows, err := db.Query(`SELECT DISTINCT code_branch FROM material_criteria WHERE status = 1 AND code_branch IS NOT NULL AND code_branch <> '' ORDER BY code_branch`).FetchAll()
	if err != nil {
		return nil, err
	}
	out := make([]Option, 0, len(rows))
	for _, row := range rows {
		v := toString(row["code_branch"])
		out = append(out, Option{Value: v, Label: v})
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
