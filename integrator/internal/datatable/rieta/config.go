package rieta

type Config struct {
	Table          string
	Joins          []string
	Columns        map[string]string
	Search         map[string][]string
	Filters        map[string]FilterConfig
	Sort           map[string]string
	PrimaryKey     string
	Where          string
	DefaultFilters []Filter
}

type FilterConfig struct {
	Field string
	Type  string
}
