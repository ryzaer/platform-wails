package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

var runtimeDLL uintptr

func NewApp() *App {
	return &App{}
}

// startup is called when the Wails application starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// // start runtime awal disini
	// go a.ensureServer()
}

// onSecondInstanceLaunch is called when the user
// tries to start Elsana while another instance is running.
func (a *App) onSecondInstanceLaunch(data options.SecondInstanceData) {
	if a.ctx == nil {
		return
	}

	// Restore window if minimized.
	runtime.WindowUnminimise(a.ctx)

	// Show existing window.
	runtime.Show(a.ctx)
}

// Open browser in native.
func (a *App) OpenBrowser(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// // ensureServer loads runtime.dll and waits
// // until the Go server is ready.
// func (a *App) ensureServer() {

// 	// Runtime already running?
// 	if a.checkHealth() {
// 		return
// 	}

// 	// Start Go runtime DLL.
// 	if err := a.startRuntime(); err != nil {

// 		runtime.MessageDialog(
// 			a.ctx,
// 			runtime.MessageDialogOptions{
// 				Title:   "Elsana",
// 				Message: err.Error(),
// 			},
// 		)

// 		return
// 	}

// 	// Wait until Go server is ready.
// 	for i := 0; i < 30; i++ {

// 		if a.checkHealth() {
// 			return
// 		}

// 		time.Sleep(500 * time.Millisecond)
// 	}

// 	runtime.MessageDialog(
// 		a.ctx,
// 		runtime.MessageDialogOptions{
// 			Title:   "Elsana",
// 			Message: "Runtime berhasil dijalankan, tetapi server tidak merespons.",
// 		},
// 	)
// }

// checkHealth checks the Go runtime.
func (a *App) checkHealth() bool {

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(
		"http://127.0.0.1:8082/api/health",
	)

	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// startRuntime loads runtime.dll and calls StartRuntime().
func (a *App) startRuntime() error {

	// --------------------------------------------------------
	// Find launcher.exe directory
	// --------------------------------------------------------

	exe, err := os.Executable()

	if err != nil {
		return fmt.Errorf(
			"gagal mendapatkan lokasi launcher: %w",
			err,
		)
	}

	appDir := filepath.Dir(exe)

	dllPath := filepath.Join(
		appDir,
		"runtime.dll",
	)

	// --------------------------------------------------------
	// Check runtime.dll
	// --------------------------------------------------------

	if _, err := os.Stat(dllPath); err != nil {

		return fmt.Errorf(
			"runtime.dll tidak ditemukan:\n\n%s",
			dllPath,
		)
	}

	// --------------------------------------------------------
	// Load runtime.dll
	// --------------------------------------------------------

	kernel32 := syscall.NewLazyDLL(
		"kernel32.dll",
	)

	loadLibrary := kernel32.NewProc(
		"LoadLibraryW",
	)

	pathPtr, err :=
		syscall.UTF16PtrFromString(
			dllPath,
		)

	if err != nil {
		return fmt.Errorf(
			"gagal membuat path DLL: %w",
			err,
		)
	}

	handle, _, _ :=
		loadLibrary.Call(
			uintptr(
				unsafe.Pointer(pathPtr),
			),
		)

	if handle == 0 {

		errCode := syscall.GetLastError()

		return fmt.Errorf(
			"gagal memuat runtime.dll\n\nPath:\n%s\n\nWindows error: %v",
			dllPath,
			errCode,
		)
	}

	runtimeDLL = handle

	// --------------------------------------------------------
	// Find StartRuntime
	// --------------------------------------------------------

	getProcAddress :=
		kernel32.NewProc(
			"GetProcAddress",
		)

	name :=
		[]byte("StartRuntime\x00")

	proc, _, _ :=
		getProcAddress.Call(
			handle,
			uintptr(
				unsafe.Pointer(
					&name[0],
				),
			),
		)

	if proc == 0 {

		return fmt.Errorf(
			"fungsi StartRuntime tidak ditemukan di runtime.dll",
		)
	}

	// --------------------------------------------------------
	// Start Go runtime
	// --------------------------------------------------------

	syscall.SyscallN(proc)

	return nil
}

// Greet returns a greeting.
func (a *App) Greet(name string) string {
	return fmt.Sprintf(
		"Hello %s, It's show time!",
		name,
	)
}
