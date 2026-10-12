// Package db is the runtime side of the case databases: the category catalogue,
// the split constants and the read-only opener. The databases are opened here
// and never created; creation, DDL and the split materialization live in the
// generation-only package internal/db/build. The read-only cross-file integrity
// check (CheckCategoriaIntegrity) lives here because the runtime runs it at startup.
// Everything is deterministic: no timestamps, fixed PRAGMAs.
package db

// Category is one of the five case labels. IDs are fixed so that corpus.db is
// byte-for-byte reproducible.
type Category struct {
	ID     int
	Nombre string
}

// Categories lists the closed set of case labels, in the order of the baseline CHECK constraint.
var Categories = []Category{
	{1, "Falta de Servicio"},
	{2, "Cobros Indebidos"},
	{3, "Problemas de Facturación"},
	{4, "Quejas"},
	{5, "Reconocimientos"},
}
