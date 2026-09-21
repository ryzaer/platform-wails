package xhttp

import (
	"app-platform/internal/config"
	"app-platform/internal/database"
	"app-platform/internal/utils"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type HandleDB = *database.Connect

func (c *Context) TokenValidity(auth ...string) (jwt.MapClaims, bool) {

	token := c.Token

	if len(auth) > 0 && auth[0] != "" {
		token = auth[0]
	}
	claims, err := utils.ReadAccessToken(token)
	if err != nil {
		return nil, true
	}

	return claims, false
}

func (c *Context) TokenExpirationStatus() (int, string) {
	// VERIFY JWT
	claims, err := c.TokenValidity()
	if err {
		return 401, "Invalid token!"
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return 401, "Invalid token date!"
	}
	// mengecek apakah ada key user
	key, ok := claims["key"].(string)
	if !ok {
		return 401, "Invalid token format!"
	}
	c.UserKey = key

	// kosongkan default token , biar bisa di set baru
	// nilai kosong menunjukkan token sudah expired, jadi harus login ulang
	c.Token = ""
	expiredAt := time.Unix(int64(exp), 0)
	// if time.Now().After(expiredAt) {
	// 	return 401, "Token expired!"
	// }

	if time.Until(expiredAt) <= config.App.TokenRefreshTime {
		// jika sudah waktunya refresh
		newToken, err := utils.GenerateAccessToken(
			config.App.TokenExpiredTime,
			map[string]any{"key": c.UserKey},
		)
		if err != nil {
			return 500, "Failed generate token"
		}
		c.Token = newToken
	}
	return 0, ""
}

func (c *Context) TokenSetSession(byOrigin ...bool) (int, string) {
	enable := false
	if len(byOrigin) > 0 {
		enable = byOrigin[0]
	}

	if enable && !c.TokenAllowOrigin() {
		return 403, "Forbidden Origin!"
	}

	if c.Token == "" {
		return 401, "No token found!"
	}

	if c.IsBearer {
		c.SetHeader("X-"+strings.ToUpper(config.App.TokenName), c.Token)
	} else {
		c.SetCookie(
			config.App.TokenName,
			c.Token,
			config.App.TokenExpiredTime,
			true,
		)
	}

	return 0, ""
}

func (c *Context) TokenGetExpiredDate(token ...string) (string, bool) {

	claims, err := c.TokenValidity(token...)
	if err {
		return "", true
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", false
	}

	return time.Unix(int64(exp), 0).Format("2006-01-02 15:04:05"), true
}

func (c *Context) TokenAllowOrigin() bool {

	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		return false
	}

	for _, o := range strings.Split(config.App.Domain, ",") {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}

	return false
}

// nanti bakal dipakai
func (c *Context) TokenGenerateByID() ([]int, string, string) {
	db, err := database.New()
	if err != nil {
		return nil, "", err.Error()
	}

	defer db.Close()

	chk_user := `
SELECT
    u.code,
    u.code_branch,
    u.role_id,
	u.name,
	u.password,
	r.scope,
	r.name AS login_as,
    r.permission_items AS permissions,
    b.secret_key
FROM users u
LEFT JOIN roles r
    ON r.id = u.role_id
LEFT JOIN branches b
    ON b.code = u.code_branch
WHERE u.code = :code
LIMIT 1
	`
	rows, err := db.Query(chk_user).Fetch(
		"code", c.UserKey,
	)
	if err != nil {
		return nil, "", err.Error()
	}

	// tentukan scope
	perm, _ := rows["permissions"].(string)
	var permissions []string
	if perm != "" {
		permissions = strings.Split(perm, ",")
	}
	if rows["scope"] == "all" {
		permissions = []string{"*"}
	}
	udata, err := json.Marshal(map[string]any{
		"code":        rows["code"],
		"name":        rows["name"],
		"scope":       rows["scope"],
		"role_id":     rows["role_id"],
		"login_as":    rows["login_as"],
		"tenant_code": rows["code_branch"],
		"permissions": permissions,
	})
	if err != nil {
		return nil, "", err.Error()
	}
	gkey, _ := utils.GenerateCode(17)
	akey := utils.EnSodiumKeyU8Api(gkey)

	sdata, err := utils.Encrypt(string(udata), gkey)
	if err != nil {
		return nil, "", err.Error()
	}

	return akey, sdata, ""
}
