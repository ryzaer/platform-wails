package socket

import (
	"log"
	"time"

	"app-platform/internal/database"
	"app-platform/internal/xhttp"
)

func handleLogin(
	ctx *xhttp.SocketContext,
	msg map[string]any,
) {

	code, ok := msg["code"].(string)
	if !ok {
		ctx.Close()
		return
	}

	db, err := database.New()
	if err != nil {
		log.Println("SOCKET DB:", err)
		ctx.Close()
		return
	}
	defer db.Close()

	sql := `
		SELECT
			name
		FROM users
		WHERE code=:code
		LIMIT 1
	`
	rows, err := db.Query(sql).Fetch("code", code)

	if err != nil {
		log.Println("SOCKET LOGIN:", err)
		ctx.Close()
		return
	}

	name, _ := rows["name"].(string)

	xhttp.Login(ctx.Client.Conn, code)

	log.Println("ARRIVE:", code)

	xhttp.Broadcast(map[string]any{
		"event": "refresh",
		"code":  code,
		"text":  "User " + code + " : " + name + " telah login!",
		"time":  time.Now().Format("2006-01-02 15:04:05"),
	})

}
