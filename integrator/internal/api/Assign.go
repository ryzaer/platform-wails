package api

import (
	"app-platform/internal/config"
	"app-platform/internal/database"
	"app-platform/internal/utils"
	"app-platform/internal/xhttp"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// fungsi check token wajib miliki middleware ElsanaAssign, karena fungsi ini untuk refresh token
func Assign(ctx *xhttp.Context) {
	// jika ada user valid & token hasil generate dari middleware/ElsanaAssign
	if ctx.Token == "" && ctx.UserKey == "" {
		ctx.Error(http.StatusUnauthorized, "Token Not Found!")
		return
	}
	// jika memang ada perbaharuan token
	if ctx.Token != "" && ctx.UserKey != "" {
		akey, adata, err := ctx.TokenGenerateByID()
		if err == "" {
			ctx.Success(map[string]any{
				"uint8a": akey,
				"assign": adata,
			})
			return
		}
	}
	ctx.Success("Assigned & applied")
}

func AssignHealth(ctx *xhttp.Context) {
	db, err := database.New()
	stdb := "On"
	if err != nil {
		stdb = "Off"
	}
	defer db.Close()
	ctx.Success(map[string]any{
		"platform": config.App.TokenName,
		"port":     config.App.Port,
		"service":  "OnDev",
		"database": stdb,
		"version":  config.App.Version,
	})
}

// ini tanpa middleware, karena ini route untuk login/logout, jadi tidak perlu token
func AssignLogin(ctx *xhttp.Context) {

	type LoginRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		CaptchaAnswer   int    `json:"captcha_answer"`
		CaptchaExpected int    `json:"captcha_expected"`
	}

	var req LoginRequest

	err := ctx.JSONBody(&req)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Data Invalid!")
		return
	}

	if req.CaptchaAnswer != req.CaptchaExpected {
		ctx.Error(http.StatusNotAcceptable, "Captcha Salah!")
		return
	}

	db, err := database.New()
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
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
WHERE u.email = :email
LIMIT 1
	`
	rows, err := db.Query(chk_user).Fetch(
		"email", req.Email,
	)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	//
	// PASSWORD VERIFY
	//
	err = bcrypt.CompareHashAndPassword(
		[]byte(rows["password"].(string)),
		[]byte(req.Password),
	)

	if err != nil {
		ctx.Error(http.StatusNotAcceptable, "Email/Password Salah!")
		return
	}

	// JWT Generate TOKEN untuk 7 hari
	token, err := utils.GenerateAccessToken(
		config.App.TokenExpiredTime,
		map[string]any{
			"key": rows["code"],
		},
	)
	if err != nil {
		ctx.Error(http.StatusNotAcceptable, "Generate Session Gagal!")
		return
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
		ctx.Error(http.StatusNotAcceptable, "JSON Token Error!")
		return
	}
	// Generate Sodium Key untuk enkripsi data ke client
	gkey, _ := utils.GenerateCode(17)
	// Generate Sodium Encrypted Data
	sdata, err := utils.Encrypt(string(udata), gkey)
	if err != nil {
		ctx.Error(http.StatusNotAcceptable, "Generate Token Gagal!")
		return
	}
	// jika semua kelar, set header dan cookie bersamaan dengan token baru
	ctx.SetHeader("X-"+strings.ToUpper(config.App.TokenName), token)
	ctx.SetCookie(
		config.App.TokenName,
		token,
		config.App.TokenExpiredTime,
		true,
	)

	// generate format Uint8Array sodium key
	akey := utils.EnSodiumKeyU8Api(gkey)
	ctx.Success(map[string]any{
		"uint8a": akey,
		"assign": sdata,
	})

}

// ini tanpa middleware, karena ini route untuk login/logout, jadi tidak perlu token
func AssignLogout(ctx *xhttp.Context) {
	// hapus token jika mode cookie,
	// mode authorization bearer tidak perlu akses route ini
	// karena client yang pegang token
	if ctx.Cookie(config.App.TokenName) != "" {
		ctx.DeleteCookie(config.App.TokenName)
		ctx.Success("Logout Success")
		return
	}
	ctx.Success("Already Logout")
}
