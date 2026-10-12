# ADR-005: despliegue del servicio web (solo loopback, un binario, sin CI ni Wine)

**Estado:** aceptada (Fase 4, Bloque D). Es la decisión que el Trabajo Escrito (TE) §6.4.1 promete documentar ("decisiones explícitamente documentadas en el ADR del despliegue").
**Implementación:** `cmd/autoatend` (flags y señales), `internal/web` (arranque, endpoints, puerto, navegador, `--selftest`), `internal/assets` (artefactos embebidos), `internal/db.CheckCategoriaIntegrity`.
**Pruebas:** CPA01–CPA05 (`internal/web/api_test.go`, `boot_test.go`, `cmd/autoatend/main_test.go`), códigos de salida de `--selftest` (`internal/web/selftest_test.go`), P99 de `/api/clasificar` (`TestLatencia_P99DeClasificar`), `go list -deps` (`TestNoDependeDeDBBuild`).

## Contexto

El servicio se presenta en la Mac propia del sustentante (macOS ARM64, Apple Silicon), que el CENEVAL autorizó (decisión del 12-oct, `PLAN.md` §5.1), sin instalación, sin permisos de administrador y sin red. Windows x86-64 (`autoatend.exe`) queda como respaldo opcional. El escenario es una demostración local ante el sínodo, no un servicio de producción (TE §6.4.1).

## Decisión

### 1. Solo loopback, sin TLS, autenticación ni CORS

El servidor escucha únicamente en `127.0.0.1`, puerto 8471 con fallback a 8472–8480 (`web.Listen`). Un socket en loopback no es accesible desde la red local ni desde el exterior, así que no hay reglas de firewall que requieran administrador ni el diálogo de permiso de red de Windows. Por la misma razón no se implementan TLS, autenticación ni CORS: no hay un canal que proteger fuera de la propia máquina, y los comentarios nunca salen de ella. La página incluye una política `Content-Security-Policy` que solo permite hablar con el mismo origen (`connect-src 'self'`), de modo que la interfaz no puede cargar ni enviar nada a otro sitio.

Riesgo aceptado: cualquier proceso o página abierta en la misma máquina puede llamar a la API. En una demostración local el único dato sensible es el corpus sintético.

### 2. Un único binario estático, sin CGo

`CGO_ENABLED=0` con `modernc.org/sqlite` (Go puro, ADR-001). `make cross` produce `dist/autoatend-macos-arm64` (Mac con Apple Silicon, principal), `dist/autoatend-macos-intel` (Mac con procesador Intel) y `dist/autoatend.exe` (Windows x86-64, respaldo) desde cualquier host, sin toolchain de C. Las compilaciones usan `-trimpath` y `-buildvcs=false`, así que el mismo código fuente produce binarios idénticos byte a byte (`make cross-sha` compila dos veces, la segunda desde una copia fuera del repositorio, y exige el mismo SHA-256 en cada binario). Los artefactos viajan dentro del ejecutable con `go:embed` (TE §7): `corpus.db`, `resultados.db`, `modelo.json`, `modelo_alternativo.json`, `gold5.json` y las páginas de `ui/`.

`go:embed` no admite `..`, así que `make assets` copia los cinco archivos de `data/` a `internal/assets/data/`. Esas copias **no se editan a mano y no se versionan** (`.gitignore`): son subproductos de `data/`, la única fuente de verdad, y versionarlas duplicaría 15 MB de binarios en cada cambio del corpus o del modelo. Todo objetivo de `make` que compila Go (`build`, `vet`, `test`, `test-integration`, `cover-gate`, `cross`) depende de `assets`, que se refresca solo si `data/` es más nuevo. En un clon limpio, `go build` sin `make assets` falla con "no matching files found" para `data/...`. La fuente de la interfaz sí está versionada: `internal/assets/ui/index.html`.

### 3. Extracción a un directorio temporal y apertura de solo lectura

Al arrancar, `web.Boot` escribe los artefactos en un directorio propio del proceso bajo el temporal del usuario (`os.MkdirTemp`, sin permisos de administrador) y abre las dos bases con `db.OpenReadOnly`, es decir `mode=ro&immutable=1` (TE §5.1.4) y URI `file:///C:/…` en Windows (`db.FileURI`). No se usa `Deserialize` (decisión ya aprobada). Con `immutable=1` SQLite no crea `-wal`, `-shm` ni archivos de bloqueo, y el servicio no puede modificar las bases: la prueba CPA03 compara el SHA-256 de `resultados.db` y `corpus.db` antes y después de inyecciones SQL.

La integridad cruzada de `categoria_id` (TE §6.1.3) se verifica en cada arranque con `ATTACH` de `corpus.db` a la conexión de `resultados.db`. Esa verificación es de solo lectura y vive en `internal/db` (`CheckCategoriaIntegrity`); `internal/db/build` delega en ella. Así `cmd/autoatend` no depende de `internal/db/build` (que crea archivos, DDL y `VACUUM`), lo que comprueba `go list -deps`. Mover la función no cambia el esquema ni el SHA del corpus.

Cualquier error de arranque (artefacto ilegible, modelo inválido, `categoria_id` huérfano, ningún puerto libre) es fatal: el programa imprime un diagnóstico en `stderr`, termina con código 1 y borra el directorio temporal. Al salir con Ctrl+C (SIGINT) o SIGTERM, el servidor se detiene y el directorio también se borra. Una terminación abrupta (por ejemplo, cerrar la ventana de consola en Windows) puede dejar el directorio `autoatend-*` en el temporal del usuario; es inocuo y el sistema lo limpia.

La consola solo emite ASCII: la página de códigos de `cmd.exe` no es UTF-8 y los acentos saldrían deformados. El JSON y el HTML sí son UTF-8.

### 4. Puerto con fallback y banner

`web.Listen` intenta 8471 y, si está ocupado, 8472 hasta 8480. El banner imprime la URL efectiva (`http://127.0.0.1:<puerto>`) y `/api/health` informa el mismo puerto (CPA05).

### 5. Apertura del navegador

Tras imprimir el banner se abre la URL en el navegador predeterminado. Es el único uso permitido de `runtime.GOOS` (`web.openBrowser`): `rundll32 url.dll,FileProtocolHandler` en Windows, `open` en macOS y `xdg-open` en Linux. Si falla, el servicio sigue y el banner pide abrir la URL a mano. `-no-browser` lo desactiva (pruebas, ejecución sin pantalla).

### 6. Interfaz de línea de comandos

`autoatend [-modelo principal|alternativo] [-no-browser] [--selftest]`. `-modelo alternativo` sirve `modelo_alternativo.json` (Naive Bayes), para la réplica de TE §6.4.2. La versión se fija al compilar con `-ldflags "-X main.version=…"` (`make` usa `VERSION`, 1.0.0 por omisión) y se reporta en `/api/health`.

### 7. Desviación aprobada: sin CI ni Wine (TE §6.3.9, Capas 1 y 2)

El TE §6.3.9 propone tres capas: CI con `windows-latest` (Capa 1), Wine en Linux (Capa 2) y `--selftest` el día del examen (Capa 3). El usuario aprobó el 11-oct (`PLAN.md` §1.3) **omitir las Capas 1 y 2** y conservar la Capa 3:

- En lugar de un runner remoto o un emulador, el binario se ejecuta de forma **nativa en la Mac (macOS ARM64)**, el equipo de la demo: `--selftest` con salida 0 y el guion de 5 pasos, verificados por el usuario el 12-oct. La ejecución en Windows queda como respaldo opcional, no verificada. Es evidencia más fuerte que la de un emulador porque es el sistema operativo real del examen.
- La principal fuente de incompatibilidades que Wine buscaría (DLLs y API Win32 nativas) no existe: sin CGo el binario es estático y solo usa la biblioteca estándar de Go y Go puro.
- La compilación cruzada (`make cross`, parte de `make verify`) cubre el "no se introduce ninguna regresión que impida la compilación" de la Capa 1.
- La Capa 3 se mantiene tal cual: `--selftest` se ejecuta en la máquina del examen antes de la presentación.

Qué verifica `--selftest` (TE §6.3.9, §6.4.3; salida 0/1, una línea por verificación, las fallidas en `stderr`; todas corren aunque una falle):

1. `corpus.db`: 73,000 comentarios, 5 categorías y ningún `categoria_id` huérfano en `comentarios`.
2. Integridad: `PRAGMA quick_check` de las dos bases y la integridad cruzada por `ATTACH`.
3. Catálogo: cada ER de `resultados.db` compila con `internal/regex` y acepta sus variantes (I-1); el vocabulario del modelo está contenido en el catálogo.
4. `gold5`: los cinco comentarios de oro embebidos se clasifican con la categoría esperada.
5. Puerto: 8471 libre o el fallback asigna otro.

## Páginas de interfaz

Además de `GET /` (la interfaz del TE §6.4.1), el servicio sirve tres páginas HTML en la misma pestaña: `/catalogo` (vista legible del catálogo), `/api` (documentación de los 5 endpoints) y `/casos-prueba` (matriz del Anexo II con su resultado medido). Son vistas sobre los mismos datos: `/catalogo` y `/api` leen los endpoints del TE, y `/casos-prueba` recibe su JSON inyectado al arrancar. **No agregan rutas de API**: el contrato sigue siendo de 5 endpoints. Cumplen las mismas reglas: CSP de mismo origen, sin recursos externos y sin `innerHTML`.

## Consecuencias

- Mediciones en Linux (`PLAN.md` §5.2): el `--selftest` tarda unas decenas de milisegundos y la latencia P99 de `/api/clasificar` está muy por debajo de 50 ms. En la Mac de la demo, `--selftest` dio 5/5; la latencia P99 no se ha medido ahí, pero el margen es de más de 400 veces (0.12 ms contra 50 ms).
- El tamaño del binario es de unos 27 MB, casi todo el corpus embebido (13.6 MB) y el motor SQLite.
- Si cambian `data/*` (otro corpus o modelo), basta `make assets` y recompilar; el binario recoge lo que `make assets` copie.
- El directorio temporal puede quedar huérfano tras una terminación abrupta (ver arriba).
