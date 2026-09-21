package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"app-platform/internal/datatable/rieta"
)

func CSV(rows []map[string]any, cfg rieta.Config) ([]byte, error) {
	return renderToBuffer(func(buf *bytes.Buffer) error {
		w := csv.NewWriter(buf)
		columns := columnNames(cfg)
		if err := w.Write(columns); err != nil {
			return err
		}
		for _, row := range rows {
			record := make([]string, len(columns))
			for i, name := range columns {
				record[i] = csvValue(row[name])
			}
			if err := w.Write(record); err != nil {
				return err
			}
		}
		w.Flush()
		return w.Error()
	})
}

func csvValue(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case fmt.Stringer:
		return v.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}
