package socket

import "app-platform/internal/xhttp"

func handleTrigger(
	ctx *xhttp.SocketContext,
	msg map[string]any,
) {

	module, ok := msg["module"].(string)
	if !ok {
		return
	}

	switch module {

	case "Inventory":
		xhttp.Broadcast(map[string]any{
			"event":  "trigger",
			"module": "Inventory",
		})

	case "chart.inserted":
		xhttp.Broadcast(map[string]any{
			"event":  "trigger",
			"module": "chart.update",
		})

	}

}
