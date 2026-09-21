export {}

declare global {
    interface Window {
        go: {
            main: {
                App: {
                    Greet(arg1: string): string
                    OpenBrowser(arg1: string): void
                }
            }
        }
    }
}