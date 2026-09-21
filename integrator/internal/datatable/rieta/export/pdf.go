package export

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"app-platform/internal/datatable/rieta"
)

// PDF writes a small dependency-free PDF table. It is deliberately conservative:
// Helvetica, landscape A4, repeated header and simple page breaks.
func PDF(rows []map[string]any, cfg rieta.Config) ([]byte, error) {
	columns := columnNames(cfg)
	if len(columns) == 0 {
		return nil, fmt.Errorf("no export columns configured")
	}

	const pageW, pageH = 842, 595
	const margin = 24
	const fontSize = 7
	const rowH = 13.0
	maxCols := 8
	if len(columns) > maxCols {
		columns = columns[:maxCols]
	}

	colW := float64(pageW-2*margin) / float64(len(columns))
	linesPerPage := int((int(pageH) - 2*margin - 28) / rowH)
	if linesPerPage < 1 {
		linesPerPage = 1
	}

	var pages []string
	for start := 0; start < len(rows) || (len(rows) == 0 && start == 0); {
		var content strings.Builder
		y := float64(pageH - margin)
		writeText := func(x, y float64, text string, bold bool) {
			font := "F1"
			if bold {
				font = "F2"
			}
			fmt.Fprintf(&content, "BT /%s %d Tf %.2f %.2f Td (%s) Tj ET\n", font, fontSize, x, y, pdfEscape(text))
		}
		for i, col := range columns {
			writeText(float64(margin)+float64(i)*colW+2, y, truncate(col, int(colW/4.2)), true)
		}
		y -= rowH
		count := 0
		for start < len(rows) && count < linesPerPage-1 {
			row := rows[start]
			for i, col := range columns {
				writeText(float64(margin)+float64(i)*colW+2, y, truncate(csvValue(row[col]), int(colW/4.2)), false)
			}
			y -= rowH
			start++
			count++
		}
		pages = append(pages, content.String())
		if len(rows) == 0 {
			break
		}
	}

	return buildPDF(pages), nil
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r", " "), "\n", " ")
	if max < 4 {
		max = 4
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-3]) + "..."
}

func pdfEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}

func buildPDF(pages []string) []byte {
	// Objects: catalog, pages, font, bold font, then page/content pairs.
	objects := make([]string, 0, 4+len(pages)*2)
	objects = append(objects,
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [PLACEHOLDER] /Count "+strconv.Itoa(len(pages))+" >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>",
	)
	kids := make([]string, 0, len(pages))
	for i, page := range pages {
		pageObj := len(objects) + 1
		contentObj := pageObj + 1
		kids = append(kids, fmt.Sprintf("%d 0 R", pageObj))
		objects = append(objects,
			fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 842 595] /Resources << /Font << /F1 3 0 R /F2 4 0 R >> >> /Contents %d 0 R >>", contentObj),
			fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(page), page),
		)
		_ = i
	}
	objects[1] = fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), len(pages))

	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n%\xE2\xE3\xCF\xD3\n")
	offsets := make([]int, len(objects)+1)
	for i, obj := range objects {
		offsets[i+1] = out.Len()
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", i+1, obj)
	}
	xref := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for i := 1; i < len(offsets); i++ {
		fmt.Fprintf(&out, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xref)
	return out.Bytes()
}
