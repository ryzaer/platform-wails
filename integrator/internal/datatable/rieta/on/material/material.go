package material

import (
	"app-platform/internal/database"
	"app-platform/internal/datatable/rieta"
	"app-platform/internal/xhttp"
)

func Handler(ctx *xhttp.Context, db *database.Connect) {
	rieta.Handler(ctx, db, Config)
}
