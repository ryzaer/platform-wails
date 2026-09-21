# README

## About

This is the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config


## Live Development

```
go install github.com/wailsapp/wails/v2/cmd/wails@latest atau
go install github.com/wailsapp/wails/v2/cmd/wails@v2.15.0 (versi 2 stabil)
wails init -l (to se js dev framework)
wails init -n YourProject -t react-ts
cd YourProject
```
To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## Custom Declare Function Bridge
Solusi yang paling aman: buat declaration sendiri
Make a file: `frontend/src/wails.d.ts`
```
export {}

declare global {
    interface Window {
        go: {
            main: {
                App: {
                    Greet(arg1: string): string
                    //...tambahkan disini
                }
            }
        }
    }
}
```
Dengan begitu TypeScript tahu struktur:
```
window
└── go
    └── main
        └── App
            ├── Greet()
            └── OpenBrowser()
```

Dan file generated Wails tetap tidak disentuh.
