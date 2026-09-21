package xhttp

import (
	"app-platform/internal/config"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)
type Middleware func(*Context) (int, string)

func Use(fn func(*Context) (int, string)) Middleware {
	return Middleware(fn)
}

func Route(route string, next HandlerFunc, origins ...any) {

	part := strings.SplitN(route, " ", 2)
	methods := strings.Split(part[0], "|")
	path := part[1]

	// if len(methods) > 0 && methods[0] == "WS" {
	// 	WebSocket(route, next func(*Context))
	// 	return
	// }

	// Default CORS
	origin := config.App.Domain
	middlewares := []Middleware{}
	for _, arg := range origins {
		switch v := arg.(type) {
		case string:
			// origin
			origin = v

		case Middleware:
			// middleware
			middlewares = append(middlewares, v)
		}

	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		ctx := New(w, r)
		corsOrigin := origin
		// PRIORITAS 1 : Authorization
		auth := r.Header.Get("Authorization")
		ctx.IsBearer = strings.HasPrefix(auth, "Bearer ")
		if ctx.IsBearer {
			ctx.Token = strings.TrimPrefix(auth, "Bearer ")
			// CORS JWT untuk request ini saja
			corsOrigin = "*"
		} else {
			// PRIORITAS 2 : Cookie
			ctx.Token = ctx.Cookie(config.App.TokenName)
		}

		// -----------------------------
		// IDENTIFIKASI CORS
		// -----------------------------
		switch corsOrigin {

		case "*":
			w.Header().Set("Access-Control-Allow-Origin", "*")

		default:
			reqOrigin := r.Header.Get("Origin")
			for _, o := range strings.Split(corsOrigin, ",") {
				o = strings.TrimSpace(o)
				if reqOrigin == o {
					w.Header().Set("Access-Control-Allow-Origin", o)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					break
				}
			}

		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-"+strings.ToUpper(config.App.TokenName))
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(append(methods, "OPTIONS"), ","))

		// -----------------------------
		// OPTIONS
		// -----------------------------

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// -----------------------------
		// METHOD
		// -----------------------------
		for _, method := range methods {

			if r.Method == method {
				// --------------------------------------
				// MIDDLEWARE AUTH SETELAH CONFIRM METHOD
				// --------------------------------------
				for _, mw := range middlewares {
					code, msgs := mw(ctx)
					if code >= 400 && code <= 599 {
						// BERHENTI JIKA MIDDLEWARE RETURN/OUTPUT CODE ERROR
						if msgs != "" {
							ctx.Error(code, msgs)
						}
						return
					}
				}
				// --------------------------------------
				// FUNCTION CONTENT
				// --------------------------------------
				next(ctx)
				return
			}

		}
		// SEND INFORMASI METHOD YG DIIZINKAN JIKA ERROR
		w.Header().Set("Allow", strings.Join(methods, ","))
		http.Error(w, `{"success":false,"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)

	})

}
