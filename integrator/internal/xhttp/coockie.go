package xhttp

import (
	"app-platform/internal/config"
	"net/http"
	"strings"
	"time"
)

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Expires  time.Time
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

// Untuk secure atau insecure IsHTTPS() di reverse proxy nginx
// location / {
//     proxy_pass http://127.0.0.1:8082;
//     proxy_set_header Host $host;
//     proxy_set_header X-Real-IP $remote_addr;
//     proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
//     proxy_set_header X-Forwarded-Proto $scheme;
// }

func (c *Context) IsHTTPS() bool {

	// HTTPS langsung
	if c.Request.TLS != nil {
		return true
	}

	// Reverse proxy (Nginx, Traefik, dll)
	if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}

	// Fallback untuk development / konfigurasi
	return strings.Contains(config.App.Domain, "https://")
}

// Membaca cookie.
func (c *Context) Cookie(name string) string {

	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// Menyimpan cookie.
func (c *Context) setCookieDuration(cookie Cookie) {

	if cookie.Path == "" {
		cookie.Path = "/"
	}

	if cookie.SameSite == 0 {
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Expires:  cookie.Expires,
		HttpOnly: cookie.HttpOnly,
		Secure:   c.IsHTTPS(),
		SameSite: cookie.SameSite,
	})
}

// Menyimpan cookie berdasarkan jumlah hari.
func (c *Context) SetCookieDays(name, value string, days int, httpOnly bool) {

	c.setCookieDuration(Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   days * 24 * 60 * 60,
		HttpOnly: httpOnly,
	})
}

// Menyimpan cookie berdasarkan jumlah tahun.
func (c *Context) SetCookieYears(name, value string, years int, httpOnly bool) {

	c.setCookieDuration(Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   years * 365 * 24 * 60 * 60,
		HttpOnly: httpOnly,
	})
}
func (c *Context) SetCookie(name, value string, duration time.Duration, httpOnly bool) {

	c.setCookieDuration(Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(duration.Seconds()),
		HttpOnly: httpOnly,
	})

}

// Menghapus cookie.
func (c *Context) DeleteCookie(name string) {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func (c *Context) HasCookie(name string) bool {
	_, err := c.Request.Cookie(name)
	return err == nil
}

// func (c *Context) RefreshCookieIfExpired(name string, db ...bool) {
// 	type Session struct {
// 		UserCode  string
// 		LoginCode string
// 		Expired   int64
// 	}

// 	// --------------------------------
// 	// Cookie ada?
// 	// --------------------------------

// 	if !c.HasCookie(name) {
// 		c.Error(http.StatusUnauthorized, "relogin")
// 		return
// 	}

// 	// --------------------------------
// 	// Decrypt Cookie
// 	// --------------------------------

// 	plain, err := utils.Decrypt(c.Cookie(name))
// 	if err != nil {
// 		c.Error(http.StatusUnauthorized, "relogin")
// 		return
// 	}

// 	// --------------------------------
// 	// Parse expired time
// 	// --------------------------------
// 	var session []any
// 	json.Unmarshal([]byte(plain), &session)
// 	userCode := session[0].(string)
// 	loginCode := session[1].(string)
// 	expired := time.Unix(
// 		int64(session[2].(float64)),
// 		0,
// 	)
// 	remain := expired.Sub(time.Now())
// 	// Sudah lewat?
// 	if remain <= 0 {
// 		c.Error(http.StatusUnauthorized, "relogin")
// 		return
// 	}
// 	// Masuk jendela refresh?
// 	if remain <= config.App.TokenRefreshTime {

// 		session[2] = time.Now().
// 			Add(config.App.TokenExpiredTime).
// 			Unix()

// 		// Encrypt
// 		// SetCookie()
// 	}

// 	// --------------------------------
// 	// Optional check database
// 	// --------------------------------
// 	if len(db) > 0 && db[0] {
// 		db, err := database.New()
// 		if err != nil {
// 			c.Error(500, err.Error())
// 			return
// 		}

// 		defer db.Close()

// 		sql := `
// 		SELECT
// 			code,
// 			token_key
// 		FROM users
// 		WHERE code=:code && token_key=:token_key
// 		LIMIT 1
// 		`
// 		rows, err := db.Query(sql).Bind(
// 			"code", userCode,
// 			"token_key", loginCode,
// 		).FetchOne()

// 		if err != nil {
// 			c.Error(http.StatusUnauthorized, "relogin")
// 			return
// 		}

// 		sql = `
// 		UPDATE users
// 		SET token_key=:token_key
// 		WHERE code=:code
// 		`
// 		newkey, _ := utils.GenerateCode(7)
// 		err = db.Query(sql).Exec(
// 			"code", rows["code"],
// 			"token_key", newkey,
// 		)
// 	}

// 	newsession, err := json.Marshal(&session)
// 	if err != nil {
// 		c.Error(http.StatusUnauthorized, "relogin")
// 		return
// 	}
// 	newCookie, _ := utils.Encrypt(newsession)
// 	//
// 	c.setCookieDuration(Cookie{
// 		Name:  name,
// 		Value: newCookie,
// 	})

// 	c.Success("resign")
// }
