package routes

import (
	"app-platform/internal/api"
	"app-platform/internal/xhttp"
)

func Stock() {
	xhttp.Route("POST /api/stock", api.Stock)
	xhttp.Route("POST /api/stock/create", api.StockCreate)
	xhttp.Route("POST /api/stock/update", api.StockUpdate)
	xhttp.Route("POST /api/stock/delete", api.StockDelete)
	xhttp.Route("POST /api/stock/restore", api.StockRestore)
	xhttp.Route("POST /api/stock/options", api.StockOptions)

	xhttp.Route("POST /api/stock/size-options", api.StockSizeOptions)
	xhttp.Route("POST /api/stock/branch-options", api.StockBranchOptions)
	xhttp.Route("POST /api/stock/export", api.StockExport)
}
