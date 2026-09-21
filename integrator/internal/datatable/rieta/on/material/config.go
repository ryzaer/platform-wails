package material

import "app-platform/internal/datatable/rieta"

var Config = rieta.Config{
	PrimaryKey: "m.id",
	Where:      "",
	Table:      "material m",
	Joins: []string{
		"JOIN suppliers s ON s.id = m.id_supplier",
		"JOIN users u ON u.code = m.created_by",
	},
	Columns: map[string]string{
		// Internal fields — dibutuhkan CRUD.
		"id":          "m.id",
		"id_supplier": "m.id_supplier",
		"code_qr":     "m.code_qr",

		// Display fields.
		"item_name":     "m.item_name",
		"supplier_name": "s.name",
		"created_at":    "m.created_at",
		"created_by":    "u.name",

		"status": "m.status",
	},
	DefaultFilters: []rieta.Filter{
		{
			Field: "status",
			Type:  "DATA_SCOPE",
			Value: "active",
		},
	},
	Search: map[string][]string{
		"global": {
			"m.item_name",
			"m.code_qr",
			"s.name",
			"u.name",
		},
		"item_name":     {"m.item_name"},
		"supplier_name": {"s.name"},
		"created_by":    {"u.name"},
		"code_qr":       {"m.code_qr"},
	},
	Filters: map[string]rieta.FilterConfig{
		"item_name": {
			Field: "m.item_name",
			Type:  "TEXT",
		},
		"supplier_name": {
			Field: "s.id",
			Type:  "MULTI_SELECTION",
		},
		"created_at": {
			Field: "m.created_at",
			Type:  "DATE_RANGE",
		},
		"created_by": {
			Field: "u.name",
			Type:  "TEXT",
		},
		"status": {
			Field: "m.status",
			Type:  "DATA_SCOPE",
		},
	},
	Sort: map[string]string{
		"item_name":     "m.item_name",
		"supplier_name": "s.name",
		"created_at":    "m.created_at",
		"created_by":    "u.name",
	},
}
