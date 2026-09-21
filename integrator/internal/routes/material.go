package routes

import (
	"app-platform/internal/api"
	"app-platform/internal/xhttp"
)

func Material() {
	xhttp.Route("POST /api/material", api.Material)
	xhttp.Route("POST /api/material/create", api.MaterialCreate)
	xhttp.Route("POST /api/material/update", api.MaterialUpdate)
	xhttp.Route("POST /api/material/delete", api.MaterialDelete)
	xhttp.Route("POST /api/material/restore", api.MaterialRestore)
	xhttp.Route("POST /api/material/options", api.MaterialOptions)
	xhttp.Route("POST /api/material/export", api.MaterialExport)
}
