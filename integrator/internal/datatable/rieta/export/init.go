package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"app-platform/internal/database"
	"app-platform/internal/datatable/rieta"
)

type Format string

const (
	FormatCSV  Format = "csv"
	FormatXLSX Format = "xlsx"
	FormatPDF  Format = "pdf"
	FormatJSON Format = "json"
)

type Request struct {
	Format string        `json:"format"`
	Query  rieta.Request `json:"query"`
}

type File struct {
	Data        []byte
	Filename    string
	ContentType string
}

func Generate(db *database.Connect, cfg rieta.Config, req Request) (File, error) {
	format := Format(strings.ToLower(strings.TrimSpace(req.Format)))
	if format == "" {
		return File{}, errors.New("export format is required")
	}

	result, err := rieta.ExecuteAll(db, cfg, req.Query)
	if err != nil {
		return File{}, err
	}

	var data []byte
	var contentType string
	var ext string

	switch format {
	case FormatCSV:
		data, err = CSV(result.Data, cfg)
		contentType, ext = "text/csv; charset=utf-8", "csv"
	case FormatXLSX:
		data, err = XLSX(result.Data, cfg)
		contentType, ext = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "xlsx"
	case FormatPDF:
		data, err = PDF(result.Data, cfg)
		contentType, ext = "application/pdf", "pdf"
	case FormatJSON:
		data, err = json.MarshalIndent(result.Data, "", "  ")
		contentType, ext = "application/json; charset=utf-8", "json"
	default:
		return File{}, fmt.Errorf("unsupported export format: %s", req.Format)
	}
	if err != nil {
		return File{}, err
	}

	return File{
		Data:        data,
		Filename:    "datatable." + ext,
		ContentType: contentType,
	}, nil
}

func renderToBuffer(fn func(*bytes.Buffer) error) ([]byte, error) {
	var buf bytes.Buffer
	if err := fn(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func columnNames(cfg rieta.Config) []string {
	keys := make([]string, 0, len(cfg.Columns))
	for key := range cfg.Columns {
		keys = append(keys, key)
	}
	// Stable ordering without importing sort into every renderer.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
