package routes

import (
	"app-platform/internal/api"
	"app-platform/internal/xhttp"
)

func Supplier() {
	xhttp.Route("POST /api/supplier", api.Supplier)
	xhttp.Route("POST /api/supplier/create", api.SupplierCreate)
	xhttp.Route("POST /api/supplier/update", api.SupplierUpdate)
	xhttp.Route("POST /api/supplier/delete", api.SupplierDelete)
	xhttp.Route("POST /api/supplier/restore", api.SupplierRestore)
	xhttp.Route("POST /api/supplier/export", api.SupplierExport)
}
