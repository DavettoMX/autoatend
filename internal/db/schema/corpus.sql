CREATE TABLE IF NOT EXISTS categorias (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT    NOT NULL UNIQUE,
    CONSTRAINT chk_nombre CHECK (nombre IN (
        'Falta de Servicio',
        'Cobros Indebidos',
        'Problemas de Facturación',
        'Quejas',
        'Reconocimientos'
    ))
);

CREATE TABLE IF NOT EXISTS comentarios (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha        TEXT    NOT NULL,
    texto        TEXT    NOT NULL,
    categoria_id INTEGER NOT NULL
                         REFERENCES categorias(id)
                         ON DELETE RESTRICT,
    CONSTRAINT chk_fecha CHECK (fecha GLOB '????-??-??')
);

CREATE INDEX IF NOT EXISTS idx_comentarios_categoria ON comentarios(categoria_id);

CREATE INDEX IF NOT EXISTS idx_comentarios_fecha ON comentarios(fecha);

CREATE INDEX IF NOT EXISTS idx_comentarios_cat_fecha ON comentarios(categoria_id, fecha);

CREATE TABLE IF NOT EXISTS particiones (
    comentario_id INTEGER PRIMARY KEY
                          REFERENCES comentarios(id)
                          ON DELETE RESTRICT,
    conjunto      TEXT    NOT NULL CHECK (conjunto IN ('train', 'test'))
);
