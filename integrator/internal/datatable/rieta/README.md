# Elsana Server-side DataTable

`internal/datatable/rieta` is the server-side query engine used by the React `rieta-server-side` component.

It is intentionally separate from `internal/database`:

- `internal/database` handles the DB connection, binding, execution and fetching.
- `internal/datatable/rieta` translates Rieta's JSON query into safe, whitelisted SQL.
- `internal/datatable/rieta/on/<module>` contains module-specific configuration.

## Structure

```text
internal/datatable/rieta/
├── config.go
├── request.go
├── response.go
├── query.go
├── handler.go
├── doc.go
└── on/
    └── stock/
        ├── config.go
        └── stock.go
```

The React and Go structures intentionally mirror each other:

```text
React                                      Go
rieta-server-side/                         internal/datatable/rieta/
└── on/                                    └── on/
    ├── stock/                                 ├── stock/
    ├── invoice/                               ├── invoice/
    └── supplier/                              └── supplier/
```

## Rieta request

```json
{
  "page": 1,
  "limit": 10,
  "search": "arabica",
  "search_by": "global",
  "sort": [
    { "field": "final_price", "direction": "desc" }
  ],
  "filters": [
    {
      "field": "code_branch",
      "type": "MULTI_SELECTION",
      "value": ["BJM"]
    }
  ]
}
```

## Stock endpoint

The API handler is `api.Material` and delegates to `internal/datatable/rieta/on/stock`.

Register it with the existing router:

```go
xhttp.Route("POST /api/material", api.Material)
```

The existing React `secureFetch()` remains the only fetch/auth layer. The Rieta datatable package does not implement authentication or HTTP fetching.

## Stock CRUD / option endpoints

The Stock module also exposes option queries for the React dialog:

- material options: `stock.MaterialOptions`
- size options: `stock.SizeOptions`
- branch options: `stock.BranchOptions`

API wrappers are in `internal/api/Material.go`.

Routes to add in `main.go`:

```go
xhttp.Route("POST /api/material/create", api.MaterialCreate)
xhttp.Route("POST /api/material/update", api.MaterialUpdate)
xhttp.Route("POST /api/material/delete", api.MaterialDelete)
xhttp.Route("POST /api/material/restore", api.MaterialRestore)
xhttp.Route("POST /api/material/options", api.MaterialOptions)
xhttp.Route("POST /api/material/size-options", api.MaterialSizeOptions)
xhttp.Route("POST /api/material/branch-options", api.MaterialBranchOptions)
```

## Export

Privileged export support lives under `internal/datatable/rieta/export`:

```text
export/
├── init.go
├── csv.go
├── xlsx.go
├── pdf.go
└── doc.go
```

The exporter receives the same `rieta.Config` and `rieta.Request` used by the
normal table query. `ExecuteAll` intentionally bypasses normal page/limit
pagination while retaining search, filters, data scope, and sorting.

Module-specific adapters belong under `on/<module>`; for example Stock exposes
an `Export` function that supplies `stock.Config` to the generic exporter.

HTTP authentication and binary response writing remain application concerns.
Rieta does not duplicate the application's auth/fetch layer.
