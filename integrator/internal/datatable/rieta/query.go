package rieta

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"app-platform/internal/database"
)

type builder struct {
	cfg    Config
	req    Request
	args   []any
	params []string
	seq    int
}

// ExecuteAll executes the same Rieta query without applying page/limit pagination.
// It is intended for privileged export operations and therefore bypasses the normal
// list-page cap while retaining the exact same filters, search and sorting.
func ExecuteAll(db *database.Connect, cfg Config, req Request) (Result, error) {
	if db == nil {
		return Result{}, errors.New("database connection is nil")
	}
	req.Normalize()
	b := &builder{cfg: cfg, req: req}

	selectSQL := b.selectSQL()
	whereSQL, err := b.whereSQL()
	if err != nil {
		return Result{}, err
	}

	countSQL := "SELECT COUNT(*) " + b.fromSQL() + whereSQL
	total, err := b.fetchCount(db, countSQL)
	if err != nil {
		return Result{}, err
	}

	dataSQL := selectSQL + " " + b.fromSQL() + whereSQL + b.orderSQL()
	rows, err := db.Query(dataSQL).FetchAll(bindArgs(b.params, b.args)...)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Data: rows,
		Pagination: Pagination{
			Page:  1,
			Limit: total,
			Total: total,
			TotalPages: func() int {
				if total > 0 {
					return 1
				}
				return 0
			}(),
		},
	}, nil
}

func Execute(db *database.Connect, cfg Config, req Request) (Result, error) {
	if db == nil {
		return Result{}, errors.New("database connection is nil")
	}
	req.Normalize()
	b := &builder{cfg: cfg, req: req}

	selectSQL := b.selectSQL()
	whereSQL, err := b.whereSQL()
	if err != nil {
		return Result{}, err
	}

	countSQL := "SELECT COUNT(*) " + b.fromSQL() + whereSQL
	total, err := b.fetchCount(db, countSQL)
	if err != nil {
		return Result{}, err
	}

	offset := (req.Page - 1) * req.Limit
	dataSQL := selectSQL + " " + b.fromSQL() + whereSQL + b.orderSQL() +
		" LIMIT :dt_limit OFFSET :dt_offset"

	dataArgs := append([]any{}, b.args...)
	dataParams := append([]string{}, b.params...)
	dataParams = append(dataParams, "dt_limit", "dt_offset")
	dataArgs = append(dataArgs, req.Limit, offset)

	rows, err := db.Query(dataSQL).FetchAll(bindArgs(dataParams, dataArgs)...)
	if err != nil {
		return Result{}, err
	}

	pages := 0
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(req.Limit)))
	}

	return Result{
		Data: rows,
		Pagination: Pagination{
			Page: req.Page, Limit: req.Limit, Total: total, TotalPages: pages,
		},
	}, nil
}

func (b *builder) selectSQL() string {
	fields := make([]string, 0, len(b.cfg.Columns))
	for key, field := range b.cfg.Columns {
		fields = append(fields, field+" AS `"+key+"`")
	}
	// Map iteration order is irrelevant to SQL semantics, but deterministic output is useful.
	// Sort by alias for stable SQL and tests.
	sortStrings(fields)
	return "SELECT " + strings.Join(fields, ", ")
}

func (b *builder) fromSQL() string {
	out := "FROM " + b.cfg.Table
	for _, join := range b.cfg.Joins {
		out += " " + join
	}
	return out
}

func (b *builder) whereSQL() (string, error) {
	clauses := make([]string, 0)
	if strings.TrimSpace(b.cfg.Where) != "" {
		clauses = append(clauses, "("+b.cfg.Where+")")
	}

	if b.req.Search != "" {
		fields := b.cfg.Search[b.req.SearchBy]
		if len(fields) == 0 {
			fields = b.cfg.Search["global"]
		}
		if len(fields) > 0 {
			search := make([]string, 0, len(fields))
			for _, field := range fields {
				name := b.param("search")
				search = append(search, field+" LIKE "+name)
				b.args = append(b.args, "%"+b.req.Search+"%")
			}
			clauses = append(clauses, "("+strings.Join(search, " OR ")+")")
		}
	}

	filters := make([]Filter, 0, len(b.cfg.DefaultFilters)+len(b.req.Filters))
	for _, defaultFilter := range b.cfg.DefaultFilters {
		filters = append(filters, defaultFilter)
	}
	for _, filter := range b.req.Filters {
		// A client-supplied filter replaces the module default for that field.
		for i := len(filters) - 1; i >= 0; i-- {
			if filters[i].Field == filter.Field {
				filters = append(filters[:i], filters[i+1:]...)
			}
		}
		filters = append(filters, filter)
	}

	for _, filter := range filters {
		cfg, ok := b.cfg.Filters[filter.Field]
		if !ok {
			continue
		}
		clause, err := b.filterSQL(cfg, filter)
		if err != nil {
			return "", err
		}
		if clause != "" {
			clauses = append(clauses, clause)
		}
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(clauses, " AND "), nil
}

func (b *builder) filterSQL(cfg FilterConfig, f Filter) (string, error) {
	switch strings.ToUpper(cfg.Type) {
	case "DATA_SCOPE":
		value, ok := f.Value.(string)
		if !ok || value == "" || strings.EqualFold(value, "all") {
			return "", nil
		}
		switch strings.ToLower(value) {
		case "active":
			p := b.param("scope")
			b.args = append(b.args, 1)
			return cfg.Field + " = " + p, nil
		case "deleted":
			p := b.param("scope")
			b.args = append(b.args, 0)
			return cfg.Field + " = " + p, nil
		default:
			return "", fmt.Errorf("unsupported datatable data scope: %s", value)
		}

	case "TEXT", "PARTIAL_MATCH":
		value, ok := f.Value.(string)
		if !ok || strings.TrimSpace(value) == "" {
			return "", nil
		}
		p := b.param("filter")
		b.args = append(b.args, "%"+strings.TrimSpace(value)+"%")
		return cfg.Field + " LIKE " + p, nil

	case "SELECTION", "MULTI_SELECTION":
		values := toSlice(f.Value)
		if len(values) == 0 {
			return "", nil
		}
		parts := make([]string, 0, len(values))
		for _, value := range values {
			p := b.param("filter")
			b.args = append(b.args, value)
			parts = append(parts, p)
		}
		return cfg.Field + " IN (" + strings.Join(parts, ", ") + ")", nil

	case "NUMBER_RANGE", "DATE_RANGE":
		if cfg.Type == "DATE_RANGE" {
			if value, ok := f.Value.(string); ok && strings.TrimSpace(value) != "" {
				p := b.param("date")
				b.args = append(b.args, strings.TrimSpace(value)[:minInt(len(strings.TrimSpace(value)), 10)])
				return "DATE(" + cfg.Field + ") = " + p, nil
			}
		}
		obj, ok := f.Value.(map[string]any)
		if !ok {
			return "", nil
		}
		parts := make([]string, 0, 2)
		if value, exists := obj["min"]; exists {
			p := b.param("min")
			b.args = append(b.args, value)
			parts = append(parts, cfg.Field+" >= "+p)
		}
		if value, exists := obj["max"]; exists {
			p := b.param("max")
			b.args = append(b.args, value)
			parts = append(parts, cfg.Field+" <= "+p)
		}
		if value, exists := obj["from"]; exists {
			p := b.param("from")
			b.args = append(b.args, value)
			parts = append(parts, cfg.Field+" >= "+p)
		}
		if value, exists := obj["to"]; exists {
			p := b.param("to")
			b.args = append(b.args, value)
			parts = append(parts, cfg.Field+" <= "+p)
		}
		return strings.Join(parts, " AND "), nil

	default:
		return "", fmt.Errorf("unsupported datatable filter type: %s", cfg.Type)
	}
}

func (b *builder) orderSQL() string {
	if len(b.req.Sort) == 0 {
		if b.cfg.PrimaryKey != "" {
			return " ORDER BY " + b.cfg.PrimaryKey + " DESC"
		}
		return ""
	}
	parts := make([]string, 0, len(b.req.Sort))
	for _, item := range b.req.Sort {
		field, ok := b.cfg.Sort[item.Field]
		if !ok {
			continue
		}
		direction := "ASC"
		if item.Direction == "desc" {
			direction = "DESC"
		}
		parts = append(parts, field+" "+direction)
	}
	if len(parts) == 0 {
		return ""
	}
	return " ORDER BY " + strings.Join(parts, ", ")
}

func (b *builder) param(prefix string) string {
	b.seq++
	name := "dt_" + prefix + "_" + strconv.Itoa(b.seq)
	b.params = append(b.params, name)
	return ":" + name
}

func (b *builder) fetchCount(db *database.Connect, sqlText string) (int, error) {
	values := bindArgs(b.params, b.args)
	row, err := db.Query(sqlText).Fetch(values...)
	if err != nil {
		return 0, err
	}
	return numberToInt(row["COUNT(*)"]), nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bindArgs(names []string, values []any) []any {
	out := make([]any, 0, len(names)*2)
	for i, name := range names {
		out = append(out, name, values[i])
	}
	return out
}

func numberToInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float64:
		return int(v)
	case []byte:
		n, _ := strconv.Atoi(string(v))
		return n
	case string:
		n, _ := strconv.Atoi(v)
		return n
	default:
		return 0
	}
}

func toSlice(value any) []any {
	switch v := value.(type) {
	case []any:
		return v
	case []string:
		out := make([]any, len(v))
		for i := range v {
			out[i] = v[i]
		}
		return out
	default:
		if value == nil {
			return nil
		}
		return []any{value}
	}
}

func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
