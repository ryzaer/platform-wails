package rieta

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type Result struct {
	Data       []map[string]any `json:"data"`
	Pagination Pagination       `json:"pagination"`
	TotalAll   int              `json:"total_all,omitempty"`
}
