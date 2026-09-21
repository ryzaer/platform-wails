import { useEffect, useState } from 'react'
import {
    Quit,
    WindowMinimise,
    WindowToggleMaximise,
} from '../wailsjs/runtime/runtime'
import { OpenBrowser } from '../wailsjs/go/main/App'
import logo from "./assets/elsana.webp";

interface HealthResult {
    database: string
    platform: string
    port: string
    service: string
    version: string
}

interface HealthResponse {
    result?: HealthResult
    success?: boolean
}

function App() {
    const [isWails, setIsWails] = useState(false)

    const [status, setStatus] = useState<
        'checking' | 'online' | 'offline'
    >('checking')

    const [health, setHealth] = useState<HealthResult | null>(null)

    const checkServer = async () => {
        setStatus('checking')

        try {
            const response = await fetch(
                'http://localhost:8082/api/health'
            )

            if (!response.ok) {
                throw new Error('Server error')
            }

            const data: HealthResponse = await response.json()

            if (!data.success || !data.result) {
                throw new Error('Invalid health response')
            }

            setHealth(data.result)
            setStatus('online')
        } catch {
            setHealth(null)
            setStatus('offline')
        }
    }

    const openApplication = () => {
        window.open('http://localhost:8082/', '_blank')
    }

    useEffect(() => {
        const runtime = Reflect.get(window, 'runtime')
        setIsWails(runtime !== undefined)

        checkServer()
    }, [])

    return (
        <div className="app">

            <div className="titlebar">
                <div className="titlebar-drag">
                    <div className="logo">
                        {/* Elsana Integrator */}
                    </div>
                </div>

                {isWails && (
                    <div className="window-buttons">
                        <button onClick={WindowMinimise}>
                            −
                        </button>

                        {/* <button onClick={WindowToggleMaximise}>
                            □
                        </button> */}

                        <button
                            className="close"
                            onClick={Quit}
                        >
                            ×
                        </button>
                    </div>
                )}
            </div>

            <main className="content">

                <div className="header">
                    <div className="brand">
                        <div className="brand-logo-box">
                            <img
                                className="brand-logo"
                                src={logo}
                                alt="Elsana"
                            />
                        </div>

                        <div className="brand-text">
                            <h1>Elsana</h1>
                            <p>App Integrator</p>
                        </div>
                    </div>
                </div>

                <section className="server-card">
                    <div className="card-header">
                        <div>
                            <h2>Server</h2>
                            <div className={`status ${status}`}>
                                <span className="status-dot" />

                                {status === 'checking' && 'Checking...'}
                                {status === 'online' && 'Online'}
                                {status === 'offline' && 'Offline'}
                            </div>
                        </div>
                    </div>

                    {health && (
                        <div className="server-info">

                            <div className="info-row">
                                <span>Platform</span>
                                <strong>
                                    {health.platform}
                                </strong>
                            </div>

                            <div className="info-row">
                                <span>Version</span>
                                <strong>
                                    {health.version}
                                </strong>
                            </div>

                            <div className="info-row">
                                <span>Port</span>
                                <strong>
                                    {health.port}
                                </strong>
                            </div>

                            <div className="info-row">
                                <span>Service</span>
                                <strong className="ok">
                                    {health.service.toUpperCase()}
                                </strong>
                            </div>

                            <div className="info-row">
                                <span>Database</span>
                                <strong className="ok">
                                    {health.database.toUpperCase()}
                                </strong>
                            </div>

                        </div>
                    )}

                    {status === 'offline' && (
                        <div className="offline-message">
                            Server tidak dapat dihubungi.
                        </div>
                    )}

                    <div className="actions">

                        <button
                            className="button secondary"
                            onClick={checkServer}
                            disabled={status === 'checking'}
                        >
                            {status === 'checking'
                                ? 'Checking...'
                                : 'Check Server'}
                        </button>

                        {/* <button
                            className="button primary"
                            onClick={openApplication}
                            disabled={status !== 'online'}
                        >
                            Open Application
                        </button> */}
                        <button className="button primary" onClick={() => OpenBrowser('http://localhost:8082')}>
                            Buka Aplikasi
                        </button>
                    </div>

                </section>

            </main>

        </div>
    )
}

export default App