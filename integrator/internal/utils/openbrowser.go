package utils

import (
	"app-platform/internal/config"
	"net"
	"os/exec"
	"runtime"
	"time"
)

// // cara pakai
// import "app-platform/internal/utils"
// if len(os.Args) > 1 && os.Args[1] != "dev" {
// 	  utils.OpenBrowser("http://localhost:" + config.App.Port)
// }

func OpenBrowser(url string) error {
	if config.App.Port != "" {
		waitServer(config.App.Port)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)

	case "darwin":
		cmd = exec.Command("open", url)

	default: // Linux
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

func waitServer(port string) {
	for {
		conn, err := net.DialTimeout(
			"tcp",
			"localhost:"+port,
			100*time.Millisecond,
		)

		if err == nil {
			conn.Close()
			return
		}

		time.Sleep(50 * time.Millisecond)
	}
}
