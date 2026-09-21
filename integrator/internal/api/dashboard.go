package api

import (
	"app-platform/internal/database"
	"app-platform/internal/xhttp"
)

func Dashboard(ctx *xhttp.Context) {

	db, err := database.New()

	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	defer db.Close()

	sql := `
SELECT
	code,
	branch,
	status
FROM branches
`

	rows, err := db.Query(sql).FetchAll()

	if err != nil {
		ctx.Error(404, "Data Tenant tidak ditemukan")
		return
	}

	ctx.Success(rows)
}
