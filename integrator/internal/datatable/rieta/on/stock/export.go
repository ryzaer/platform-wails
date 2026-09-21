package stock

import (
	"app-platform/internal/database"
	rietaexport "app-platform/internal/datatable/rieta/export"
)

// Export generates a Stock export using the same Rieta configuration used by
// the server-side table. HTTP authentication/response handling stays outside
// this package.
func Export(db *database.Connect, req rietaexport.Request) (rietaexport.File, error) {
	return rietaexport.Generate(db, Config, req)
}
