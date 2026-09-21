# platform-wails

A prototype **Go + Wails application framework** for building desktop and server-capable applications with a shared application layer.

This project is intentionally **not tied to Elsana**.

It is a framework/prototype environment for testing an architecture where the same Go application can provide:

- a Wails desktop application
- a browser/server application
- HTTP APIs
- WebSocket/realtime communication
- embedded React frontends
- native launcher integration
- external processing/rendering components when required

The project can be used as a foundation for different kinds of applications.

---

## Status

> **Prototype / Experimental Framework**

`platform-wails` is an architecture prototype.

Its purpose is to validate a reusable application structure before adopting the structure in a production project.

The code, directory structure, runtime commands, and integration strategy may change during development.

---

## Concept

The central idea is to separate the **application layer** from the environment in which the application runs.

The same Go application logic should be able to run through different hosts:

```text
                    Application
                         │
              ┌──────────┴──────────┐
              │                     │
         Desktop Host          Server Host
              │                     │
            Wails                HTTP Server
              │                     │
              └──────────┬──────────┘
                         │
                    Go Application
                         │
          ┌──────────────┼──────────────┐
          │              │              │
         API         WebSocket         Core
          │              │              │
          └──────────────┴──────────────┘
                         │
                    Data / Services
```

The frontend can then be shared between environments:

```text
             React Frontend
                    │
          ┌─────────┴─────────┐
          │                   │
       Wails                 Browser
          │                   │
          └─────────┬─────────┘
                    │
                Go API
                    │
                Go Core
```

The framework does not require a particular business domain.

It can be used for an inventory system, administration application, analytics platform, document system, internal tool, or another application with similar architectural requirements.

---

## Goals

The prototype is intended to explore the following:

- reusable Go application architecture
- Wails desktop integration
- React frontend integration
- HTTP API
- WebSocket communication
- shared application logic
- desktop and server execution modes
- native launcher integration
- optional external processing components
- simple inter-process communication
- clean separation between UI, application logic, and infrastructure

The objective is **not** to build another framework with unnecessary abstraction.

The structure should remain understandable and practical.

---

## Project Structure

```text
platform-wails/
│
├── integrator/
│   ├── build/
│   ├── frontend/
│   └── internal/
│       ├── api/
│       ├── config/
│       ├── database/
│       ├── datatable/
│       ├── middleware/
│       ├── routes/
│       ├── socket/
│       ├── utils/
│       └── xhttp/
│
├── launcher/
│
├── python/
│
├── app.go
├── main.go
├── go.mod
├── go.sum
├── wails.json
├── build.bat
├── .gitignore
└── README.md
```

---

## Directory Roles

### `integrator/`

The main application integration area.

It contains the frontend and application-side Go packages used by the runtime.

```text
integrator/
├── build/
├── frontend/
└── internal/
```

The name `integrator` is intentional: its responsibility is to bring the different application components together.

---

### `integrator/frontend/`

The frontend application.

The current prototype uses React, but the architecture does not depend on a particular frontend framework.

The important requirement is that the frontend can operate in both:

```text
Desktop
  ↓
Wails
```

and:

```text
Browser
  ↓
HTTP/WebSocket server
```

without duplicating the application UI.

---

### `integrator/internal/`

The application/backend layer.

Current areas include:

```text
internal/
├── api/
├── config/
├── database/
├── datatable/
├── middleware/
├── routes/
├── socket/
├── utils/
└── xhttp/
```

These packages should contain application functionality rather than Wails-specific UI logic.

The goal is to keep the core application usable even when Wails is not running.

---

### `launcher/`

Native launcher/integration code.

This prototype includes a dedicated **C++ native launcher** for Windows.

The launcher is not the application business layer. Its purpose is to act as a small native entry point and integration layer between the operating-system process environment and the application runtime.

The launcher can:

- start the Go/Wails runtime
- load native DLL components
- bridge native Windows functionality
- forward command-line arguments
- connect the Go runtime with the Python renderer through DLL-based integration

The Windows launcher uses the **Windows GUI subsystem** (`/SUBSYSTEM:WINDOWS`) so the user-facing application can run without opening a console window.

The architecture can therefore be viewed as:

```text
Windows
   │
   ▼
C++ Launcher
   │
   ├── Go / Wails DLL
   │
   └── Python / Nuitka DLL
```

The important idea is that Go and Python do not need to be exposed as separate user-facing executables for the desktop application.

The C++ layer provides the native process/DLL boundary, while Go remains responsible for the application/runtime layer and Python remains an optional processing/renderer layer.

For example:

```text
C++ Launcher
     │
     ├── integrator.dll
     │       └── Go / Wails
     │
     └── renderer.dll
             └── Python / Nuitka
```

This native launcher is particularly useful when a project requires Windows-specific integration while still keeping the main application logic in Go.

The launcher is intentionally kept separate from the application logic.

---

### `python/`

Optional processing/renderer components.

Python can be used for tasks where another runtime is more appropriate than Go, for example:

- document processing
- rendering
- image processing
- data processing
- specialized libraries

The framework does not require Python.

It is simply an optional integration point.

---

## Runtime Model

The framework is intended to support several execution modes.

Conceptually:

```text
application
│
├── desktop
│   └── Wails + frontend + Go application
│
├── server
│   └── HTTP + WebSocket + Go application
│
└── core
    └── command/processing operations
```

A runtime may expose commands such as:

```text
runtime --start gui
runtime --start api
runtime --core <command> [arguments]
```

The exact command-line interface is intentionally not fixed yet.

Different projects using this framework may expose different commands.

---

## Core Command and Serial Output

The prototype supports the idea of an optional serial output channel.

A command can run normally:

```text
runtime --core hello
```

or provide a result file:

```text
runtime --core hello --serial tmp/serial-XXXXXXX.txt
```

The important design rule is:

> `--serial` is an optional output channel, not a separate execution mode.

This allows an external requester to own the request/result lifecycle.

For example:

```text
Requester
    │
    │ generate serial
    ▼
runtime --core command --serial result-file
    │
    ▼
application/core
    │
    ▼
result-file
    │
    ▼
Requester reads result
    │
    └── removes temporary result
```

The requester could be:

- another Go process
- a desktop application
- a service
- a test program
- another application component

This keeps the runtime independent from the consumer of its output.

---

## Native Runtime Integration

One of the distinctive parts of this prototype is the use of a native C++ launcher as a bridge between multiple runtimes.

The intended Windows structure is:

```text
                 Windows Application
                        │
                        ▼
                 C++ Native Launcher
                        │
              ┌─────────┴─────────┐
              │                   │
              ▼                   ▼
        Go / Wails DLL       Python DLL
              │                   │
              ▼                   ▼
        Application          Processing /
        Runtime              Renderer
```

The C++ launcher uses native Windows facilities and DLL loading to connect these components inside the desktop application environment.

This makes the launcher a **runtime integration layer**, rather than a second application containing business logic.

The prototype explores a model where:

- C++ handles native Windows/subsystem integration
- Go handles application/runtime logic
- Wails provides the desktop application host
- React provides the frontend
- Python provides optional processing/rendering capabilities

The components communicate through explicit boundaries rather than requiring every part of the system to be implemented in the same language.

---

## Application vs Host

One of the main architectural principles is separating **what the application does** from **where it runs**.

### Application layer

Responsible for:

```text
Business logic
Database access
API handlers
WebSocket events
Core commands
Services
Configuration
```

### Desktop host

Responsible for:

```text
Wails
Native window
Desktop lifecycle
Desktop-specific integration
```

### Server host

Responsible for:

```text
HTTP server
WebSocket server
Server lifecycle
Remote access
```

The application layer should not need to know whether it is currently running inside a Wails window or a normal server process unless that distinction is explicitly required.

---

## Frontend Strategy

The frontend should communicate through application APIs rather than directly depending on the desktop host.

Conceptually:

```text
React
 │
 ├── HTTP API
 │
 └── WebSocket
       │
       ▼
     Go
```

For desktop applications, Wails-specific bindings may be added where they provide a clear advantage.

However, Wails bindings should not become the only way to access application functionality when the same functionality is expected to work in server mode.

---

## Design Principles

### 1. Simple over abstract

The framework should avoid abstraction that does not solve a real problem.

A developer should be able to open the project and understand where the application code lives.

### 2. Shared application logic

Desktop and server modes should reuse the same application logic whenever practical.

```text
             Application
             /         \
        Desktop       Server
          Wails        HTTP
```

not:

```text
Desktop Application
       +
Separate Server Application
       +
Duplicated Business Logic
```

### 3. Frontend independence

The frontend should not be tightly coupled to the Wails environment.

The same UI should be capable of running in a browser when the application is exposed as a web application.

### 4. Infrastructure stays separate

Database, network transport, desktop integration, and external processing should remain separable from business logic.

### 5. Optional components

Not every project needs every component.

For example:

```text
Project A
  Go + Wails + React

Project B
  Go + Wails + React + WebSocket

Project C
  Go + API + React + Python renderer
```

The framework should accommodate these differences without forcing unused components into every application.

### 6. Requester owns temporary results

When `--serial` is supplied, the requester owns the result file lifecycle.

This makes concurrent requests easier to manage:

```text
serial-A.txt
serial-B.txt
serial-C.txt
```

Each request can have its own correlation identifier.

---

## Example Architecture

A complete application may eventually look like:

```text
                     ┌──────────────┐
                     │ React        │
                     │ Frontend     │
                     └──────┬───────┘
                            │
                ┌───────────┴───────────┐
                │                       │
             Wails                   Browser
                │                       │
                └───────────┬───────────┘
                            │
                       Go Runtime
                            │
              ┌─────────────┼─────────────┐
              │             │             │
             API        WebSocket        Core
              │             │             │
              └─────────────┴─────────────┘
                            │
                     Application Logic
                            │
             ┌──────────────┼──────────────┐
             │              │              │
          Database       Services       External
                                         Workers
```

The external worker area is optional.

---

## Why This Prototype Exists

A production application often starts with a simple executable and gradually accumulates:

```text
GUI
API
database
WebSocket
background jobs
renderers
native integrations
```

Without a clear boundary, these components can become tightly coupled.

This prototype explores whether Wails can serve as the desktop host while the Go application remains reusable outside the desktop environment.

The goal is to validate the architecture **before** committing a production project to it.

---

## Relationship to Other Projects

`platform-wails` is a generic prototype.

A production project can adopt all or only part of its architecture.

For example:

```text
platform-wails
      │
      │ prototype
      ▼
validated architecture
      │
      ├── Project A
      ├── Project B
      └── Project C
```

The prototype should therefore avoid domain-specific names, assumptions, and business rules.

---

## Development Status

Current work focuses on validating:

- Wails integration
- Go runtime structure
- React frontend embedding
- application package organization
- native launcher integration
- core command execution
- serial-based result communication
- desktop/server separation

Features and directory names may change as the architecture evolves.

---

## Future Possibilities

Depending on the project, the architecture may later support:

- authentication
- background workers
- scheduled jobs
- plugin systems
- file processing
- document rendering
- realtime events
- desktop notifications
- local databases
- remote databases
- service discovery
- multiple frontend applications

These are possibilities, not requirements of the framework.

---

## Philosophy

`platform-wails` is intended to be a **practical application skeleton**, not a framework that hides the application behind layers of abstraction.

The preferred structure is:

```text
easy to read
      ↓
easy to modify
      ↓
easy to debug
      ↓
easy to reuse
```

The framework should make application development easier without becoming the application itself.

---

**Go + Wails Application Framework Prototype**
