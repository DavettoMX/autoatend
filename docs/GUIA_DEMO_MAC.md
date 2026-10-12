# Guía paso a paso: demo de AutoAtend en la Mac

Esta guía explica cómo preparar la Mac **una vez, en casa**, y cómo presentar la demo **el día del examen**, siguiendo el guion de 5 pasos del Trabajo Escrito (§6.4.3).

No necesitas instalar nada. AutoAtend es **un solo archivo** que ya trae dentro las bases de datos, el catálogo de expresiones regulares, el modelo y la página web. Funciona **sin internet**.

> **Convención:** los bloques grises son comandos. Se escriben en la app **Terminal**, uno a la vez, y se ejecutan con la tecla **Enter**. Lo que va después de `#` es una explicación y no hace falta escribirlo.

---

## Parte 1. Preparación en casa (una sola vez, ~15 minutos, con internet)

### 1.1 Averigua qué procesador tiene tu Mac

1. Haz clic en el menú  (arriba a la izquierda) → **Acerca de esta Mac**.
2. Mira la línea **Chip** o **Procesador**:
   - Si dice **Apple M1, M2, M3 o M4**, usarás el archivo **`autoatend-macos-arm64`**.
   - Si dice **Intel**, usarás el archivo **`autoatend-macos-intel`**.

### 1.2 Descarga el programa desde GitHub

1. Abre Safari o Chrome e inicia sesión en **github.com** con tu cuenta **DavettoMX**. El repositorio es privado, así que sin sesión no vas a ver nada.
2. Entra a: **https://github.com/DavettoMX/autoatend/releases/tag/v1.0.0**
3. En la sección **Assets**, haz clic en el archivo que te tocó en el paso 1.1 (`autoatend-macos-arm64` o `autoatend-macos-intel`) y también en **`SHA256SUMS.txt`**.
4. Los dos archivos quedan en tu carpeta **Descargas**.

### 1.3 Abre la Terminal

1. Pulsa **⌘ Cmd + Espacio**, escribe `Terminal` y pulsa **Enter**.
2. Se abre una ventana con texto. Ahí se escriben los comandos de esta guía.

### 1.4 Crea una carpeta para AutoAtend y mueve el programa ahí

Copia y pega estos comandos, uno por uno. Si tu Mac es **Intel**, cambia `arm64` por `intel` en el segundo comando.

```bash
mkdir -p ~/AutoAtend                                       # crea la carpeta AutoAtend en tu usuario
mv ~/Downloads/autoatend-macos-arm64 ~/AutoAtend/autoatend   # mueve el programa y lo renombra "autoatend"
cd ~/AutoAtend                                             # entra a la carpeta
```

> Si tu Mac está en español y no encuentra `~/Downloads`, no cambies nada: la carpeta "Descargas" se llama `Downloads` por dentro.

### 1.5 Dale permiso de ejecución y quita el bloqueo de macOS

macOS bloquea los programas que se descargan de internet y no vienen de la App Store. Estos dos comandos lo resuelven:

```bash
chmod +x autoatend                              # permiso para ejecutarlo como programa
xattr -d com.apple.quarantine autoatend          # quita la marca de "descargado de internet"
```

- Si el segundo comando responde **`No such xattr`**, no pasa nada: significa que no había bloqueo.
- Ninguno de los dos pide contraseña ni permisos de administrador.

### 1.6 Comprueba que el archivo es el correcto

```bash
shasum -a 256 autoatend
```

Debe imprimir exactamente uno de estos códigos (también vienen en `SHA256SUMS.txt`):

| Archivo | SHA-256 |
|---|---|
| `autoatend-macos-arm64` | `624f5c827dba3075aa0ac07bcf8369851a8f3caeb978d8cc052e8ba6e603f100` |
| `autoatend-macos-intel` | `e96b9370db956186d34508cc2fabbdcf72d52efda303a30100c4e520a0e1a3a7` |

Si el código no coincide, el archivo se descargó mal: bórralo y descárgalo otra vez.

### 1.7 Primera prueba: el autochequeo (`--selftest`)

```bash
./autoatend --selftest
```

Debe mostrar **cinco líneas `ok`** y terminar con `todas las verificaciones pasaron`:

```text
ok     corpus.db: 73000 comentarios en 5 categorias
ok     integridad: quick_check ok en corpus.db y resultados.db; ...
ok     catalogo: 6196 lexemas, 6196 expresiones regulares compiladas y 6891 variantes aceptadas (I-1); ...
ok     gold5: 5 de 5 comentarios de oro clasificados con su categoria esperada (modelo rls/linea_base)
ok     puerto: 127.0.0.1:8471 libre
selftest: todas las verificaciones pasaron (version 1.0.0)
```

Para confirmar el código de salida, escribe:

```bash
echo $?        # debe imprimir 0 (cero = todo bien; 1 = algo falló)
```

**Si macOS muestra el aviso "no se puede abrir porque proviene de un desarrollador no identificado":**
1. Haz clic en **Aceptar** o **Cancelar**.
2. Abre **Configuración del Sistema → Privacidad y seguridad**.
3. Baja hasta el mensaje sobre "autoatend" y haz clic en **Abrir de todos modos**. Puede pedirte tu contraseña de la Mac.
4. Vuelve a la Terminal y repite `./autoatend --selftest`. Si aparece otra ventana de aviso, elige **Abrir**.

Solo pasa la primera vez.

### 1.8 Segunda prueba: el servicio web con la interfaz

```bash
./autoatend
```

1. La Terminal muestra:
   ```text
   AutoAtend 1.0.0 - clasificador de comentarios de TelArg
   Servicio listo en http://127.0.0.1:8471
   Abra esa direccion en el navegador. Para terminar presione Ctrl+C o cierre esta ventana.
   ```
2. Tu navegador se abre solo en **http://127.0.0.1:8471** con la página de AutoAtend. Si no se abre, copia esa dirección en Safari o Chrome.
3. Escribe un comentario y pulsa **Clasificar** (o **⌘ Cmd + Enter**).
4. Para **detener** el servicio, regresa a la Terminal y pulsa **Control + C**. Debe decir `Servicio detenido.`

> **Importante:** mientras la demo esté en curso, **no cierres la ventana de Terminal**. Si la cierras, el servicio se detiene y la página deja de responder.

### 1.9 Prueba sin internet

Apaga el Wi-Fi y repite los pasos 1.7 y 1.8. Todo debe funcionar igual, porque el servicio solo usa `127.0.0.1`, que es tu propia Mac.

### 1.10 (Opcional) Abrirlo con doble clic

En el **Finder**, abre la carpeta `AutoAtend` (en tu carpeta de usuario) y haz doble clic en `autoatend`. Se abre una Terminal que ejecuta el servicio y después el navegador. Esto equivale al paso 1 del guion del TE ("ejecutar el binario con doble clic").

### 1.11 Prepara la Mac para la presentación

- Lleva el **cargador** y los **adaptadores** (USB-C a HDMI, por ejemplo). La guía del CENEVAL pide traerlos y considerar el tiempo de configuración.
- Activa **No molestar**: Centro de control → Concentración → **No molestar**.
- Evita que la pantalla se apague: Configuración del Sistema → **Pantalla bloqueada** → "Apagar la pantalla…" en **Nunca** mientras dure el examen.
- Cierra las apps que no uses y deja abiertas solo **Terminal**, el **navegador** y tu **presentación**.
- Instala las actualizaciones de macOS **antes** del día. La guía advierte que no hay tiempo para que el equipo se actualice durante el examen.
- En la Terminal, agranda la letra para que se lea en el proyector: **⌘ Cmd + +** (varias veces).

---

## Parte 2. El día del examen

### 2.1 Antes de entrar (llega 30 minutos antes, como pide el CENEVAL)

1. Conecta la Mac al proyector y comprueba que se ve.
2. Abre la Terminal y ejecuta:
   ```bash
   cd ~/AutoAtend
   ./autoatend --selftest
   ```
   Confirma las cinco líneas `ok`. Esto ya es la "Capa 3" de verificación del TE §6.3.9.
3. Deja la Terminal abierta y lista. Todavía no lances el servicio.

### 2.2 El guion de 5 pasos (TE §6.4.3)

**Paso 1. Ejecutar el binario.**
- Escribe `./autoatend` en la Terminal, o haz doble clic en el Finder.
- Muestra el **banner** con la dirección, y luego el navegador abriéndose solo.
- Qué decir: *"Es un único ejecutable, sin instalación ni dependencias. Trae dentro las dos bases de datos, el catálogo de expresiones regulares y el modelo. Escucha solo en 127.0.0.1, así que no necesita permisos de red."*

**Paso 2. Clasificar dos comentarios en vivo.**

Escríbelos a mano frente al sínodo. Estos son los resultados que se midieron con el modelo actual:

| Tipo | Comentario sugerido | Resultado esperado |
|---|---|---|
| **Claro** | `Hice mi recarga de $349.00 hace quince días y sigo sin internet.` | **Falta de Servicio**, ~99.7 % |
| **Frontera** | `Me atendieron por correo y quedó resuelto, aunque tuve que insistir 5 veces. Gracias.` | **Reconocimientos** 52 % frente a **Quejas** 47 % |
| **Frontera, con sarcasmo** (opcional, tu ejemplo) | `Es increíble la falta de atención! Mi negocio sigue sin internet después de más de una semana y diversas quejas por distintos medios.` | **Falta de Servicio** 48 % frente a **Quejas** 46 % |

- Qué decir en el caso frontera: *"Aquí el modelo duda entre dos categorías próximas. Es justo la confusión que anticipé en el diagnóstico (§5.2.2) y la razón de usar Macro-F1."*
- Si escribes un texto distinto, el resultado puede variar. Es normal.

**Paso 3. Explicar la decisión con `lexemas_clave`.**
- Señala la tabla **"Lexemas clave"** bajo el resultado. Cada fila tiene el lexema, **su expresión regular del catálogo** y cuánto aportó a la decisión.
- Qué decir: *"La decisión se explica en términos de los lexemas que generó el analizador de la consigna b. Por ejemplo, `sig(e|o)` es la expresión regular que agrupa las variantes 'sigue' y 'sigo'."*

**Paso 4. Mostrar el catálogo en `/api/lexemas` (conecta la consigna e con la b).**
- Haz clic en el enlace del pie de la página, o escribe en el navegador:
  `http://127.0.0.1:8471/api/lexemas?orden=variantes`
- Muestra 3 expresiones regulares con alternancia, por ejemplo:
  - `cobr(ado|ando|ar|aron|o|os)`
  - `pag(ado|ando|ar|e|o|os)`
  - `aclara(ción|do|ndo|r|ron)`
- Para buscar una en particular, usa `q`. Por ejemplo: `http://127.0.0.1:8471/api/lexemas?q=factur`
- Qué decir: *"El servicio consume en tiempo real el catálogo que generó el software base: forma canónica, expresión regular y frecuencias por categoría."*

> Si el navegador muestra el JSON en una sola línea difícil de leer, en Firefox y Chrome se ve formateado automáticamente. En Safari puedes activar "Vista de texto sin formato".

**Paso 5. Ejecutar `--selftest`.**
1. En la Terminal, abre una **pestaña nueva** con **⌘ Cmd + T**. No cierres la que tiene el servicio.
2. Ejecuta:
   ```bash
   cd ~/AutoAtend
   ./autoatend --selftest
   ```
3. Como el servicio ya ocupa el puerto 8471, la línea del puerto dirá `8471 ocupado; el fallback asigna 127.0.0.1:8472`. **Es correcto**: demuestra el mecanismo de respaldo de puertos.
4. Qué decir: *"Son los criterios de aceptación del plan de pruebas, verificados en vivo en menos de un segundo: integridad del corpus, I-1 sobre todo el catálogo, los 5 comentarios de oro y el puerto."*

### 2.3 Extra para la réplica: el modelo alternativo

Si te preguntan por la comparación entre Naive Bayes y Regresión Logística (TE §6.4.2):
1. Detén el servicio con **Control + C**.
2. Ejecuta:
   ```bash
   ./autoatend -modelo alternativo
   ```
3. Clasifica el mismo comentario frontera y compara las confianzas. Ahora responde el modelo Naive Bayes Multinomial.
4. Para volver al modelo principal: **Control + C** y luego `./autoatend`.

### 2.4 Al terminar

Pulsa **Control + C** en la Terminal del servicio. Al salir, el programa borra solo sus archivos temporales.

---

## Parte 3. Problemas comunes y su solución

| Síntoma | Causa | Solución |
|---|---|---|
| `zsh: permission denied: ./autoatend` | Falta el permiso de ejecución | `chmod +x autoatend` |
| `zsh: no such file or directory: ./autoatend` | No estás en la carpeta correcta | `cd ~/AutoAtend` y vuelve a intentar |
| Aviso de "desarrollador no identificado" | Bloqueo de macOS (Gatekeeper) | Paso 1.7: Configuración del Sistema → Privacidad y seguridad → **Abrir de todos modos** |
| `zsh: bad CPU type in executable` | Descargaste el archivo del procesador equivocado | Repite el paso 1.1 y descarga el otro archivo |
| El banner dice `:8472` (u otro puerto) en lugar de `:8471` | Ya hay otro AutoAtend o programa en el 8471 | Es normal (respaldo automático). Usa la dirección que muestra el banner. Para cerrar otros: `pkill autoatend` |
| El navegador no se abre solo | Bloqueo o navegador por defecto | Copia la dirección del banner en Safari o Chrome |
| La página dice que no puede conectarse | Cerraste la Terminal del servicio | Vuelve a ejecutar `./autoatend` |
| El resultado muestra el aviso "Ninguna palabra del comentario pertenece al vocabulario del modelo" | El texto solo tiene palabras vacías, números o palabras que el modelo no conoce | Escribe una oración más completa |

### Plan B, si la interfaz web no se ve en el proyector

El TE §6.4.3 prevé usar `--selftest` como evidencia: ejecútalo en la Terminal (Parte 2, paso 5). Además puedes clasificar desde la Terminal sin navegador, mientras el servicio está corriendo en otra pestaña:

```bash
curl -s -X POST http://127.0.0.1:8471/api/clasificar -H 'Content-Type: application/json' -d '{"texto":"Hice mi recarga y sigo sin internet"}'
```

### Plan C: compilar desde el código en la Mac (solo si los binarios fallan; requiere internet)

1. Instala las herramientas de desarrollo de Apple (incluyen `git` y `make`):
   ```bash
   xcode-select --install
   ```
2. Instala **Go 1.26 o superior**: descarga el instalador `.pkg` para macOS desde https://go.dev/dl/ y ábrelo.
3. Descarga el código. El repositorio es privado: inicia sesión con `gh auth login` si tienes GitHub CLI, o usa la URL HTTPS con tu usuario y un token.
   ```bash
   git clone https://github.com/DavettoMX/autoatend.git ~/autoatend-src
   cd ~/autoatend-src
   ```
4. Compila y prueba:
   ```bash
   make build                   # copia los artefactos y compila dist/autoatend
   ./dist/autoatend --selftest
   ./dist/autoatend
   ```

---

## Parte 4. Datos útiles para responder preguntas

- **Versión:** 1.0.0. `GET /api/health` responde `{"ok":true,"version":"1.0.0","lexemas":6196,"puerto":8471}`.
- **Corpus:** 73,000 comentarios, semilla 10433, SHA-256 `37f285ff…`.
- **Catálogo:** 6,196 lexemas. I-1: 0 fallos; I-2: 0 filas; I-3: SHA idéntico en dos corridas.
- **Modelo:** Regresión Logística Softmax. Macro-F1 0.9795 en test temporal; 0.835 con plantillas no vistas (cifra realista). Detalles en `docs/reports/evaluacion_modelo.md` y en ADR-004.
- **Latencia P99** de `/api/clasificar`: ~0.1–0.2 ms medida en Linux. Si quieres la cifra de tu Mac, ejecuta la prueba en el repositorio compilado (Plan C): `go test -run TestLatencia_P99DeClasificar -v ./internal/web/`.
