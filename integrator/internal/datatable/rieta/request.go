package rieta

import (
	"encoding/json"
	"strings"
)

type Request struct {
	Page     int      `json:"page"`
	Limit    int      `json:"limit"`
	Search   string   `json:"search"`
	SearchBy string   `json:"search_by"`
	Sort     []Sort   `json:"sort"`
	Filters  []Filter `json:"filters"`
	Tab      string   `json:"tab"`
}

type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type Filter struct {
	Field string `json:"field"`
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

func Parse(body []byte) (Request, error) {
	req := Request{Page: 1, Limit: 10}
	if len(strings.TrimSpace(string(body))) == 0 {
		return req, nil
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return req, err
	}
	req.Normalize()
	return req, nil
}

func (r *Request) Normalize() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 10
	}
	if r.Limit > 500 {
		r.Limit = 500
	}
	r.Search = strings.TrimSpace(r.Search)
	if r.SearchBy == "" {
		r.SearchBy = "global"
	}
	for i := range r.Filters {
		if r.Filters[i].Field == "" {
			r.Filters[i].Field = r.Filters[i].Key
		}
	}
	for i := range r.Sort {
		r.Sort[i].Direction = strings.ToLower(r.Sort[i].Direction)
		if r.Sort[i].Direction != "desc" {
			r.Sort[i].Direction = "asc"
		}
	}
}
