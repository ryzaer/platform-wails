package api

import (
	"app-platform/internal/database"
	"app-platform/internal/datatable/rieta/export"
	"app-platform/internal/datatable/rieta/on/supplier"
	"app-platform/internal/xhttp"
)

func Supplier(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()

	supplier.Handler(ctx, db)
}

func SupplierCreate(ctx *xhttp.Context) {
	var req supplier.SaveRequest

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

	id, err := supplier.Create(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{
		"id": id,
	})
}

func SupplierUpdate(ctx *xhttp.Context) {
	var req supplier.SaveRequest

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

	if err := supplier.Update(db, req); err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{
		"id": req.ID,
	})
}

func SupplierDelete(ctx *xhttp.Context) {
	var req supplier.IDsRequest

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

	affected, err := supplier.SoftDelete(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{
		"affected": affected,
	})
}

func SupplierRestore(ctx *xhttp.Context) {
	var req supplier.IDsRequest

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

	affected, err := supplier.Restore(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	ctx.Success(map[string]any{
		"affected": affected,
	})
}

func SupplierExport(ctx *xhttp.Context) {
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

	file, err := supplier.Export(db, req)
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
