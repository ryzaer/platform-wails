package socket

import "app-platform/internal/xhttp"

func ApplyHandle(ctx *xhttp.SocketContext) {

	defer ctx.Close()

	for {

		var msg map[string]any

		if err := ctx.Read(&msg); err != nil {
			break
		}

		event, ok := msg["event"].(string)
		if !ok {
			continue
		}

		switch event {

		case "login":
			handleLogin(ctx, msg)

		case "trigger":
			handleTrigger(ctx, msg)

		case "chat":
			handleChat(ctx, msg)

		}

	}

}
