package api

import (
	"app-platform/internal/database"
	"app-platform/internal/datatable/rieta/export"
	"app-platform/internal/datatable/rieta/on/material"
	"app-platform/internal/xhttp"
)

func Material(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	material.Handler(ctx, db)
}

func MaterialCreate(ctx *xhttp.Context) {
	var req material.SaveRequest

	if err := ctx.JSONBody(&req); err != nil {
		ctx.Error(400, "invalid request body")
		return
	}

	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	id, err := material.Create(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{"id": id})
}

func MaterialUpdate(ctx *xhttp.Context) {
	var req material.SaveRequest

	if err := ctx.JSONBody(&req); err != nil {
		ctx.Error(400, "invalid request body")
		return
	}

	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	if err := material.Update(db, req); err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{"id": req.ID})
}

func MaterialDelete(ctx *xhttp.Context) {
	var req material.IDsRequest

	if err := ctx.JSONBody(&req); err != nil {
		ctx.Error(400, "invalid request body")
		return
	}

	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	affected, err := material.SoftDelete(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{"affected": affected})
}

func MaterialRestore(ctx *xhttp.Context) {
	var req material.IDsRequest

	if err := ctx.JSONBody(&req); err != nil {
		ctx.Error(400, "invalid request body")
		return
	}

	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	affected, err := material.Restore(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{"affected": affected})
}

func MaterialOptions(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	result, err := material.SupplierOptions(db)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Success(result)
}

func MaterialExport(ctx *xhttp.Context) {
	var req export.Request

	if err := ctx.JSONBody(&req); err != nil {
		ctx.Error(400, "invalid request body")
		return
	}

	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	file, err := material.Export(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.SetHeader("Content-Type", file.ContentType)
	ctx.SetHeader(
		"Content-Disposition",
		`attachment; filename="`+file.Filename+`"`,
	)
	ctx.SetHeader(
		"Access-Control-Expose-Headers",
		"Content-Disposition",
	)

	ctx.Writer.WriteHeader(200)
	_, _ = ctx.Writer.Write(file.Data)
}
