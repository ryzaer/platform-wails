package stock

import "app-platform/internal/datatable/rieta"

var Config = rieta.Config{
	PrimaryKey: "mc.id",
	Where:      "",
	Table:      "material_criteria mc",
	Joins: []string{
		"JOIN material m ON m.id = mc.id_material",
		"JOIN sizes s ON s.id = mc.id_size",
	},
	Columns: map[string]string{
		"id":          "mc.id",
		"item_name":   "m.item_name",
		"id_material": "mc.id_material",
		"id_size":     "mc.id_size",
		"quantity":    "mc.quantity",
		"size_name":   "s.name",
		"size_price":  "mc.size_price",
		"final_price": "mc.final_price",
		"create_at":   "mc.create_at",
		"code_branch": "mc.code_branch",
		"status":      "mc.status",
	},
	DefaultFilters: []rieta.Filter{{Field: "status", Type: "DATA_SCOPE", Value: "active"}},
	Search: map[string][]string{
		"global": {
			"m.item_name",
			"m.code_qr",
			"s.name",
			"mc.code_branch",
		},
		"item_name": {"m.item_name"},
		"code_qr":   {"m.code_qr"},
		"size_name": {"s.name"},
	},
	Filters: map[string]rieta.FilterConfig{
		"code_branch": {Field: "mc.code_branch", Type: "MULTI_SELECTION"},
		"item_name":   {Field: "m.item_name", Type: "TEXT"},
		"size_name":   {Field: "s.name", Type: "MULTI_SELECTION"},
		"size_price":  {Field: "mc.size_price", Type: "NUMBER_RANGE"},
		"final_price": {Field: "mc.final_price", Type: "NUMBER_RANGE"},
		"quantity":    {Field: "mc.quantity", Type: "NUMBER_RANGE"},
		"create_at":   {Field: "mc.create_at", Type: "DATE_RANGE"},
		"status":      {Field: "mc.status", Type: "DATA_SCOPE"},
		"id_size":     {Field: "mc.id_size", Type: "MULTI_SELECTION"},
	},
	Sort: map[string]string{
		"item_name":   "m.item_name",
		"id_material": "mc.id_material",
		"id_size":     "mc.id_size",
		"quantity":    "mc.quantity",
		"size_name":   "s.name",
		"size_price":  "mc.size_price",
		"final_price": "mc.final_price",
		"create_at":   "mc.create_at",
		"code_branch": "mc.code_branch",
		"status":      "mc.status",
	},
}
