package api

import (
	"app-platform/internal/database"
	"app-platform/internal/datatable/rieta/export"
	"app-platform/internal/datatable/rieta/on/stock"
	"app-platform/internal/xhttp"
)

func Stock(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()
	stock.Handler(ctx, db)
}

func StockCreate(ctx *xhttp.Context) {
	var req stock.SaveRequest
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
	id, err := stock.Create(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}
	ctx.Success(map[string]any{"id": id})
}

func StockUpdate(ctx *xhttp.Context) {
	var req stock.SaveRequest
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
	if err := stock.Update(db, req); err != nil {
		ctx.Error(400, err.Error())
		return
	}
	ctx.Success(map[string]any{"id": req.ID})
}

func StockDelete(ctx *xhttp.Context) {
	var req stock.IDsRequest
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
	affected, err := stock.SoftDelete(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}
	ctx.Success(map[string]any{"affected": affected})
}

func StockRestore(ctx *xhttp.Context) {
	var req stock.IDsRequest
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
	affected, err := stock.Restore(db, req)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}
	ctx.Success(map[string]any{"affected": affected})
}

func StockOptions(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()
	result, err := stock.MaterialOptions(db)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Success(result)
}

func StockSizeOptions(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()
	result, err := stock.SizeOptions(db)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Success(result)
}

func StockBranchOptions(ctx *xhttp.Context) {
	db, err := database.New()
	if err != nil {
		ctx.Error(500, "database connection failed")
		return
	}
	defer db.Close()
	result, err := stock.BranchOptions(db)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Success(result)
}

func StockExport(ctx *xhttp.Context) {

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

	file, err := stock.Export(db, req)
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
