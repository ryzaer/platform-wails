package api

import (
	"app-platform/internal/config"
	"app-platform/internal/xhttp"
	"fmt"
)

func Ping(ctx *xhttp.Context) {

	ctx.SetCookieYears(
		config.App.TokenName,
		"12345",
		1,
		true,
	)

	fmt.Println("Cookie Header:", ctx.Request.Header.Get("Cookie"))

	for _, c := range ctx.Request.Cookies() {
		fmt.Printf("%s=%s\n", c.Name, c.Value)
	}

	rslt := "tidak ada"
	if ctx.HasCookie(config.App.TokenName) {
		rslt = "ada"
	}

	ctx.Success(rslt)
}
func PingCheckSession(ctx *xhttp.Context) {

	rslt := "Tidak ada cookie"
	if ctx.Token != "" {
		ctx.Request.Header.Get("Cookie")
		rslt = ctx.Token

		// fmt.Println("Check Cookie:", getCookie)

		// for _, c := range ctx.Request.Cookies() {
		// 	fmt.Printf("%s=%s\n", c.Name, c.Value)
		// }
	}

	ctx.Success(rslt)
}

func PingLoginSession(ctx *xhttp.Context) {

	rslt := "Session Sudah Ada"
	if !ctx.HasCookie(config.App.TokenName) {
		ctx.SetCookieYears(
			config.App.TokenName,
			"12345",
			1,
			true,
		)
		rslt = "session elsana-app:12345 dibuat"

		fmt.Println("Login Cookie:", ctx.Cookie(config.App.TokenName))
	}

	ctx.Success(rslt)
}

func PingLogoutSession(ctx *xhttp.Context) {

	rslt := "Session Sudah Dihapus"
	if ctx.HasCookie(config.App.TokenName) {
		ctx.DeleteCookie(config.App.TokenName)
		rslt = "Session elsana-app:12345 Baru Saja Dihapus"
		fmt.Println("Logout Cookie:", ctx.Cookie(config.App.TokenName))
	}

	ctx.Success(rslt)
}
