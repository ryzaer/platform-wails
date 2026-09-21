package material

import (
	"app-platform/internal/database"
	rietaexport "app-platform/internal/datatable/rieta/export"
)

func Export(db *database.Connect, req rietaexport.Request) (rietaexport.File, error) {
	return rietaexport.Generate(db, Config, req)
}
