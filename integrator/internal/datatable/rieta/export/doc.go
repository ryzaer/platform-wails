// Package export provides Rieta's privileged export subsystem.
//
// Exporters are intentionally independent from individual table modules. A
// module supplies its Rieta Config and query; this package only turns the
// resulting rows into CSV, XLSX, PDF, or JSON bytes.
package export
