package middleware

import (
	"app-platform/internal/config"
	"app-platform/internal/database"
	"app-platform/internal/xhttp"
	"net/http"
	"strings"
)

// middleware ini return harus success, jika masa expired token belum habis
// jika masa token sudah habis, maka return error 401. jika masa token
// mendekati habis, maka akan di generate token baru, memperbaharui
// cookie server serta dikirim ke client melalui
// header X-ELSANA-TOKEN jika mode bearer
func ElsanaAssign(ctx *xhttp.Context) (int, string) {

	code, err := ctx.TokenExpirationStatus()
	if code != 0 {
		return code, err
	}

	if ctx.UserKey != "" {
		db, err := database.New()
		if err != nil {
			return http.StatusServiceUnavailable, err.Error()
		}

		defer db.Close()

		chkKey := `
		SELECT code
		FROM users
		WHERE code=:code
			AND status=1
		LIMIT 1
		`

		user, err := db.Query(chkKey).Fetch(
			"code", ctx.UserKey,
		)

		if err != nil {
			return http.StatusServiceUnavailable, err.Error()
		}

		if len(user) == 0 {
			return http.StatusUnauthorized, "User not found"
		}
	}

	if ctx.Token != "" {
		// Kirim jika memang ada pembaharuan token
		ctx.SetHeader("X-"+strings.ToUpper(config.App.TokenName), ctx.Token)
		ctx.TokenSetSession(true)
	}

	return 0, ""

}

func ElsanaSecure(ctx *xhttp.Context) (int, string) {
	_, err := ctx.TokenExpirationStatus()
	if err == "" {
		// Kirim jika memang ada pembaharuan token
		if ctx.Token != "" {
			ctx.SetHeader("X-"+strings.ToUpper(config.App.TokenName), ctx.Token)
		}
		return 0, ""
	}

	return http.StatusUnauthorized, err
}

func Elsana(ctx *xhttp.Context) (int, string) {
	_, err := ctx.TokenExpirationStatus()
	if err == "" {
		// harus return 0, ""
		// true wajib dengan domain terdaftar
		return ctx.TokenSetSession(true)
	}
	return http.StatusUnauthorized, err
}
