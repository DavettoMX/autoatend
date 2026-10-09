// Package db owns the SQLite schemas (corpus.db, resultados.db), their fixed
// build-time PRAGMAs, the cross-file integrity helper and the persisted
// train/test split. Everything here is deterministic: no timestamps, fixed PRAGMAs.
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

// corpusDDL is applied in order inside one transaction (baseline Anexo I, part A).
// Deviation (ADR-007): the journal mode is DELETE, not WAL; it is set through the
// connection PRAGMAs in open.go.
var corpusDDL = []string{
	`CREATE TABLE categorias (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT    NOT NULL UNIQUE,
    CONSTRAINT chk_nombre CHECK (nombre IN (
        'Falta de Servicio',
        'Cobros Indebidos',
        'Problemas de Facturación',
        'Quejas',
        'Reconocimientos'
    ))
)`,
	`CREATE TABLE comentarios (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha        TEXT    NOT NULL,
    texto        TEXT    NOT NULL,
    categoria_id INTEGER NOT NULL
                         REFERENCES categorias(id)
                         ON DELETE RESTRICT,
    CONSTRAINT chk_fecha CHECK (fecha GLOB '????-??-??')
)`,
	`CREATE INDEX idx_comentarios_categoria ON comentarios(categoria_id)`,
	`CREATE INDEX idx_comentarios_fecha ON comentarios(fecha)`,
	`CREATE INDEX idx_comentarios_cat_fecha ON comentarios(categoria_id, fecha)`,
	// Phase 1 addition: the materialized deterministic train/test assignment (see split.go).
	`CREATE TABLE particiones (
    comentario_id INTEGER PRIMARY KEY
                          REFERENCES comentarios(id)
                          ON DELETE RESTRICT,
    conjunto      TEXT    NOT NULL CHECK (conjunto IN ('train', 'test'))
)`,
}

// resultadosDDL is applied in order inside one transaction (baseline Anexo I, part B).
// Deviation (ADR-006): uq_frecuencias_global restores uniqueness of the global
// (categoria_id IS NULL) row, which the composite PK cannot enforce because
// SQLite treats NULLs as distinct in a PRIMARY KEY.
var resultadosDDL = []string{
	`CREATE TABLE corridas (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    semilla       INTEGER NOT NULL,
    parametros    TEXT    NOT NULL,
    total_lexemas INTEGER NOT NULL CHECK(total_lexemas > 0)
)`,
	`CREATE TABLE lexemas (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    forma_canonica    TEXT    NOT NULL UNIQUE,
    expresion_regular TEXT    NOT NULL,
    corrida_id        INTEGER NOT NULL
                              REFERENCES corridas(id)
)`,
	// categoria_id references corpus.db's categorias(id); that cross-file check
	// lives in attach.go (generation pipeline only, ADR-007).
	`CREATE TABLE frecuencias (
    lexema_id    INTEGER NOT NULL REFERENCES lexemas(id),
    corrida_id   INTEGER NOT NULL REFERENCES corridas(id),
    categoria_id INTEGER,
    valor        INTEGER NOT NULL CHECK(valor > 0),
    PRIMARY KEY (lexema_id, corrida_id, categoria_id)
)`,
	`CREATE TABLE variantes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    lexema_id  INTEGER NOT NULL REFERENCES lexemas(id),
    forma      TEXT    NOT NULL,
    frecuencia INTEGER NOT NULL CHECK(frecuencia > 0),
    UNIQUE (lexema_id, forma)
)`,
	`CREATE INDEX idx_frecuencias_lexema ON frecuencias(lexema_id)`,
	`CREATE INDEX idx_frecuencias_corrida ON frecuencias(corrida_id)`,
	`CREATE INDEX idx_variantes_lexema ON variantes(lexema_id)`,
	`CREATE UNIQUE INDEX uq_frecuencias_global ON frecuencias(lexema_id, corrida_id) WHERE categoria_id IS NULL`,
}
