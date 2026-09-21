package xhttp

import (
	"embed"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
)

func Spa(frontend embed.FS) http.Handler {

	sub, err := fs.Sub(frontend, "dist")
	if err != nil {
		panic(err)
	}

	return SpaHandler(http.FS(sub))
}

func SpaDev() http.Handler {

	// Jika sudah ada React build
	// return SpaHandler(http.Dir("./cmd/elsana_new/dist"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Development Mode</title>
<style>
html,body{
	height:100%;
	margin:0;
	font-family:Arial,Helvetica,sans-serif;
	background:#fff;
}
body{
	display:flex;
	align-items:center;
	justify-content:center;
}
h1{
	font-size:36px;
	font-weight:500;
	color:#444;
}
</style>
</head>
<body>
	<h1>Page In Development</h1>
</body>
</html>`))

	})

}

func SpaHandler(fsys http.FileSystem) http.Handler {

	file := http.FileServer(fsys)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// API tidak diproses disini
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Coba buka file
		_, err := fsys.Open(strings.TrimPrefix(r.URL.Path, "/"))

		// Kalau file tidak ada, kirim index.html
		if err != nil {

			r2 := *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = "/"

			file.ServeHTTP(w, &r2)
			return
		}

		// File ditemukan
		file.ServeHTTP(w, r)

	})

}
