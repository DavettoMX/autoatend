# ADR-003: stoplist, homógrafos y criterio de agrupación de variantes

**Estado:** aceptada (Fase 2, Bloque A). Las decisiones 1 a 4 de la sección "Decisiones" las tomó el usuario el 11-oct-2026.
**Implementación:** `internal/lexer` (`normalize.go`, `tokenizer.go`, `stoplist.go` con `stoplist/v1.json`, `stem.go`, `symtab.go`, `regexgen.go`, `prefix.go`) y `internal/sweep`.
**Pruebas:** CPU09 a CPU14, CPE01 a CPE04, CPP01 (generador), CPP02 y CPP03.

## Contexto

El Trabajo Escrito (TE) §6.2.2 pide un analizador léxico que normaliza, tokeniza, filtra palabras funcionales y agrupa las variantes de superficie de cada palabra en un lexema con su expresión regular (ER) y sus frecuencias. Deja abiertos cuatro puntos, que este ADR cierra:

1. Qué palabras funcionales se filtran, dado que el propio TE nombra homógrafos (`bajo`, `para`, `sobre`, `entre`, `como`, `más`, `solo`, `aun`) que son a la vez palabra funcional y palabra de contenido.
2. Cómo se tratan los números, las variantes con ruido (`segun`, `ole`, `q`, `x`) y los acentos.
3. Qué criterio agrupa `factura`, `facturas` y `facturación` en un mismo lexema (el TE no define ninguno).
4. Qué se hace con los folios de la forma `TK-35699A`.

Cuando las listas del TE chocan entre sí, se resuelve **a favor del caso** (`2907_10433.pdf`), de acuerdo con la precedencia de fuentes del proyecto. La lista literal vive en `internal/lexer/stoplist/v1.json` (46 entradas) y `TestStoplistFileMatchesPolicy` comprueba que contiene exactamente la política siguiente, sin una entrada más ni una menos.

## Decisiones

### 1. Política de la stoplist (versión `v1`)

El archivo es una tabla cerrada `forma → clase`, versionada y embebida con `go:embed` (mismo patrón que `internal/db/schema/schema.go`). `lexer.Lookup(forma)` devuelve la clase, y `lexer.Version()` la versión que se registra en `corridas.parametros`.

| Clase | Formas |
|---|---|
| `articulo` | el, la, los, las, un, una, unos, unas, lo; contracciones **al, del** |
| `preposicion` | de, en, con, por, a, ante, tras, durante, mediante, según, sin, hasta, desde, contra, **hacia** |
| `conjuncion` | y, e, o, u, ni, que, pero, sino, aunque, porque, pues |
| `interjeccion` | ah, eh, ay, uf, olé |
| variantes de superficie | `segun` (preposición), `ole` (interjección), `q` (conjunción) y `x` (preposición) |

- Las contracciones `al` y `del` se clasifican como `articulo` para conservar las cuatro clases cerradas del caso. Es una etiqueta, no hay efecto sobre el filtrado.
- **Variantes de superficie.** `segun` y `ole` son las formas sin tilde de `según` y `olé` (la consigna e admite entrada libre). `q` y `x` son las abreviaturas que el generador de ruido escribe por `que` y por `por`: aparecen en 5,763 y 3,126 comentarios del corpus. Sin ellas, esas abreviaturas llegarían al catálogo como lexemas (hueco detectado durante el desarrollo). CPU14 exige que cada variante quede en la misma clase que su forma original.
- **Se retienen los homógrafos del TE:** `bajo`, `para`, `sobre`, `entre`, `como`, `más`, `solo`, `aun`. La tabla es plana y consulta solo la forma, así que no distingue el uso funcional del uso de contenido (`bajo` preposición contra `bajo` adjetivo; `solo` adverbio contra `sólo`). Se prefiere conservar a perder señal: el costo es que algunos usos funcionales llegan al catálogo (sub-filtrado). Ese costo es el que mide L-1 (TE §6.3.10) sobre una muestra de 200 apariciones. La limitación se declara; no se intenta resolver por contexto.
- **No se filtran** `se`, `no`, `este` ni `eso`, porque no pertenecen a ninguna de las cuatro clases (pronombre, adverbio de negación, demostrativos). Tampoco `él`, `qué`, `sí` ni `té`, que se distinguen de `el`, `que`, `si` y `te` precisamente por la tilde. CPU14 verifica todos estos casos.
- Consecuencia medida: entre los lexemas más frecuentes del catálogo figuran `me` (52,023), `mi` (42,070), `no` (34,722) y `ya`, `es`, `esto`. Son palabras funcionales que la política deja pasar a propósito, porque no pertenecen a las cuatro clases cerradas.

### 2. Sin plegado de acentos

No hay plegado global de acentos. `facturación` y `facturacion` son dos formas de superficie distintas, y también lo son `según` y `segun`, que por eso aparecen por separado en la stoplist. El plegado fusionaría pares como `él`/`el`, `qué`/`que`, `sí`/`si` y `más`/`mas`, justo los que la política necesita distinguir. El ruido de tildes del corpus (`TildeDrop` 0.80) tiene como costo lexemas duplicados sin acento, que L-2 contribuye a medir.

### 3. Números: el AFD los reconoce y no se catalogan

El tokenizador es un AFD explícito de cuatro estados (INICIO, PALABRA, NUMERO, SEPARADOR) con tabla `[estado][clase]` y tres clases de entrada (LETRA, DIGITO, SEP), fiel al TE §6.2.2. Emite `Token{Text, Kind}` con `Kind` igual a `Palabra` o `Numero`. `Analyze` descarta los `Numero`: un precio o un número de ticket no es un lexema del dominio y haría crecer el catálogo sin aportar señal (`$1,499.00` queda como `1`, `499`, `00`, los tres descartados).

Un cambio de letra a dígito, o de dígito a letra, cierra el token en curso y abre otro: `abc123` produce `abc` y `123`. Está probado en `TestTokenize_LetterDigitBoundarySplits`.

### 4. NFC a mano, restringida al alfabeto del TE

`Normalize` pasa a minúsculas y compone, con una tabla de 14 pares, una vocal más U+0301, la `u` más U+0308 y la `n` más U+0303, tanto en minúsculas como en mayúsculas, hasta el alfabeto `[a-záéíóúüñ]`. Cualquier otro carácter fuera de LETRA, DIGITO o SEPARADOR (signos, símbolos, emojis, letras de otros alfabetos, marcas combinantes sueltas) se reemplaza por un espacio. Se escribió a mano porque `golang.org/x/text` está prohibido (no se añaden dependencias). Es O(|texto|) e idempotente. No intenta ser una NFC completa: solo cubre el alfabeto del caso.

### 5. Stemmer ligero propio como criterio de agrupación

El TE §6.2.2 pide agrupar variantes en un lexema pero no define el criterio. Se usa un stemmer de sufijos del español, escrito a mano (`internal/lexer/stem.go`, versión `v1`):

- 29 sufijos ordenados de mayor a menor longitud: `aciones, iciones, amiento, imiento, ación, ición, mente, ieron, iendo, aron, ando, ados, adas, idos, idas, ado, ada, ido, ida, ar, er, ir, as, es, os, a, e, o, s`.
- Se quita **solo el primer sufijo que coincide** (el más largo), y solo si la raíz conserva al menos 3 runas; de lo contrario la forma queda intacta y no se prueba un sufijo más corto.
- Sin plegado de acentos (decisión 2).
- La clave de agrupación es la raíz. **El nombre del lexema es su variante más frecuente** (empate: la menor en orden lexicográfico), no la raíz. Por eso el catálogo muestra `factura` con la ER `factur(a|ación|as)`.

Es un criterio de agrupación, no un análisis morfológico, y sus límites son medibles. Con el corpus completo (6,257 lexemas):

- **L-2:** 5,727 de 6,257 lexemas (91.5 %) tienen una sola variante observada. La mayoría son errores tipográficos o formas raras, cada uno con su propio lexema (`cobraan`, `cobranndo`).
- **Sobreagrupación:** la raíz mínima de 3 runas funde palabras distintas: `est(a|ado|e|o)`, `pla(n|no)`, `sol(a|o)`, `hac(e|er)`.
- **Subagrupación:** `cobraron` (`cobr(ando|ar|aron|o|os|s)`) y `cobran` (`cobra(n|na|no)`) quedan en lexemas distintos, porque `n` no es un sufijo de la lista.

Se descartó el lematizador completo o un diccionario: añaden datos o dependencias y el TE no los pide. Se aceptó el stemmer ligero porque su costo es medible (L-2) y reversible: la versión se guarda en `corridas.parametros` (`stemmer: "v1"`), de modo que un cambio posterior es trazable.

### 6. Generador de ER por prefijo común

`GenerateER` ordena y deduplica las variantes por bytes, toma el prefijo común más largo (en runas) y construye `prefijo(s1|s2|...)`. La gramática del TE §6.2.3 no admite alternativas vacías, así que si una variante coincide con el prefijo, el prefijo retrocede una runa (`cobro` y `cobros` dan `cobr(o|os)`, `mes` y `meses` dan `me(s|ses)`). Escapa con `\` los metacaracteres `| * + ? ( ) [ ] \ -`. Con una sola variante la ER es la forma misma; con un prefijo vacío, una alternancia sin prefijo.

## Propuesta de mejora (no implementada): folios clave-valor

El corpus incluye folios con prefijo, guion y dígitos: 1,089 comentarios contienen `TKT-` seguido de un número (por ejemplo `TKT-583128`), y la propuesta cubre también formatos como `TK-35699A` o `REP-2026-35`. Con el AFD actual el guion es un separador, así que `TKT-583128` se parte en `tkt` (palabra, que entra al catálogo como lexema con 1,136 apariciones) y `583128` (número, descartado). El folio se pierde como unidad, y los de forma `TK-35699A` dejarían además un `a` suelto (descartado como preposición).

**Propuesta:** catalogar el folio completo como un token único. Para ello:

- se añade `-` al alfabeto de entrada (el TE §6.2.1 ya menciona los "guiones"), y
- el AFD incorpora un estado ALFANUM que absorbe secuencias mixtas de letras, dígitos y guiones internos, con una transición de salida bien definida cuando el guion no está rodeado de alfanuméricos.

Encaja en el apartado de trabajo futuro del TE §7.4. **No se implementa en este bloque**: modifica el alfabeto y el AFD entregados a los sinodales (el TE está congelado) y, de acuerdo con las reglas del proyecto, cualquier diferencia con el TE requiere la aprobación previa del usuario. Si se aprueba, el cambio toca `normalize.go` (el guion deja de ser separador), `tokenizer.go` (un estado y una columna más en `delta`), sus pruebas y la versión del tokenizador en `corridas.parametros`.

## Consecuencias

- `lexer.Analyze` es la **única definición léxica**: el barrido, el simulador y, después, la inferencia en Go y el vectorizador de Python (como subproceso) la comparten, así que catálogo y clasificador no pueden discrepar sobre qué es una palabra.
- `corridas.parametros` registra la versión de la stoplist, la del stemmer, el alfabeto, la regla de números, el corte de train, el total de comentarios y la configuración de ruido, sin timestamps ni rutas, para que dos barridos sean byte a byte idénticos (I-3).
- La stoplist es una tabla plana: no se resuelve la ambigüedad funcional/contenido de los homógrafos. El sub-filtrado se mide con L-1 y se reporta; no se corrige con reglas de contexto.
- Cambiar la política exige una nueva versión (`v2`) del archivo y de este ADR.
