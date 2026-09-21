package supplier

import "app-platform/internal/datatable/rieta"

var Config = rieta.Config{
	PrimaryKey: "s.id",
	Table:      "suppliers s",

	Joins: []string{
		"JOIN users u ON u.code = s.created_by",
	},

	Columns: map[string]string{
		"id":         "s.id",
		"name":       "s.name",
		"contact":    "s.contact",
		"email":      "s.email",
		"address":    "s.address",
		"latitude":   "s.latitude",
		"longitude":  "s.longitude",
		"notes":      "s.notes",
		"rating":     "s.rating",
		"status":     "s.status",
		"created_by": "u.name",
		"created_at": "s.created_at",
	},

	DefaultFilters: []rieta.Filter{
		{Field: "status", Type: "DATA_SCOPE", Value: "active"},
	},

	Search: map[string][]string{
		"global": {
			"s.name",
			"s.contact",
			"s.email",
			"u.name",
		},
		"name":    {"u.name"},
		"contact": {"s.contact"},
		"email":   {"s.email"},
	},

	Filters: map[string]rieta.FilterConfig{
		"rating": {
			Field: "s.rating",
			Type:  "NUMBER_RANGE",
		},
		"created_at": {
			Field: "s.created_at",
			Type:  "DATE_RANGE",
		},
	},

	Sort: map[string]string{
		"name":       "s.name",
		"contact":    "s.contact",
		"email":      "s.email",
		"rating":     "s.rating",
		"created_at": "s.created_at",
		"created_by": "u.name",
	},
}
