package rieta

import (
	"errors"
	"io"

	"app-platform/internal/database"
	"app-platform/internal/xhttp"
)

func Handler(ctx *xhttp.Context, db *database.Connect, cfg Config) {
	if ctx == nil || db == nil {
		if ctx != nil {
			ctx.Error(500, "datatable is not initialized")
		}
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Error(400, "invalid request body")
		return
	}

	req, err := Parse(body)
	if err != nil {
		ctx.Error(400, "invalid datatable query")
		return
	}

	result, err := Execute(db, cfg, req)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Success(result)
}

var ErrInvalidConfig = errors.New("invalid datatable config")
