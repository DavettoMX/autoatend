CREATE TABLE corridas (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    semilla       INTEGER NOT NULL,
    parametros    TEXT    NOT NULL,
    total_lexemas INTEGER NOT NULL CHECK(total_lexemas > 0)
);

CREATE TABLE lexemas (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    forma_canonica    TEXT    NOT NULL UNIQUE,
    expresion_regular TEXT    NOT NULL,
    corrida_id        INTEGER NOT NULL
                              REFERENCES corridas(id)
);

CREATE TABLE frecuencias (
    lexema_id    INTEGER NOT NULL REFERENCES lexemas(id),
    corrida_id   INTEGER NOT NULL REFERENCES corridas(id),
    categoria_id INTEGER,
    valor        INTEGER NOT NULL CHECK(valor > 0),
    PRIMARY KEY (lexema_id, corrida_id, categoria_id)
);

CREATE TABLE variantes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    lexema_id  INTEGER NOT NULL REFERENCES lexemas(id),
    forma      TEXT    NOT NULL,
    frecuencia INTEGER NOT NULL CHECK(frecuencia > 0),
    UNIQUE (lexema_id, forma)
);

CREATE INDEX idx_frecuencias_lexema ON frecuencias(lexema_id);

CREATE INDEX idx_frecuencias_corrida ON frecuencias(corrida_id);

CREATE INDEX idx_variantes_lexema ON variantes(lexema_id);

CREATE UNIQUE INDEX uq_frecuencias_global ON frecuencias(lexema_id, corrida_id) WHERE categoria_id IS NULL;
