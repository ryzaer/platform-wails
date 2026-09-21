package socket

import "app-platform/internal/xhttp"

func handleChat(
	ctx *xhttp.SocketContext,
	msg map[string]any,
) {

	to, ok := msg["to"].(string)
	if !ok {
		return
	}

	text, ok := msg["text"].(string)
	if !ok {
		return
	}

	_ = xhttp.Send(to, map[string]any{
		"event": "chat",
		"text":  text,
	})

}
