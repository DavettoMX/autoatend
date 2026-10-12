# Contexto de trabajo para Claude Code (traspaso a la Mac)

> Este archivo reúne lo que Claude Code tenía en su memoria local y en la conversación de desarrollo (10–12 de octubre de 2026). Sirve para que una sesión nueva, por ejemplo en la Mac, tenga el contexto completo. **Léelo junto con `CLAUDE.md` y `PLAN.md`.**

## 1. Situación

- **Qué es:** el caso práctico 10433 del CENEVAL (Acuerdo 286, Licenciatura en Ciencias Computacionales). AutoAtend clasifica comentarios de clientes de TelArg en 5 categorías.
- **Examen oral (EXOAL):** **13 de octubre de 2026, en la CDMX.** La demo se presenta en la **Mac del usuario** (MacBook Air, Apple Silicon), que el CENEVAL autorizó. Windows queda como respaldo opcional.
- **Estado:** las fases 1–4 están implementadas. La v1.0.3 está publicada y verificada en la Mac (selftest 5/5 y guion de demo).
- **Guía de la demo paso a paso:** `docs/GUIA_DEMO_MAC.md`.

## 2. Dónde está cada cosa

| Qué | Dónde |
|---|---|
| Código (solo `proyect/`) | GitHub privado **DavettoMX/autoatend** (rama `main`) |
| Binarios listos | Releases de ese repo; el último es **v1.0.3** (`autoatend-macos-arm64`, SHA `f160200e…15d0`) |
| Repo local completo (incluye PDFs) | PC de torre del usuario: `~/proyectos/ceneval-projecto-exoal` |
| Caso, TE y guías del CENEVAL (PDF) | **No están en GitHub** a propósito: el PDF del caso prohíbe su reproducción e intercambio. El usuario debe copiarlos a la Mac por USB si los necesita para la presentación |

**El repo de GitHub es una exportación** de `proyect/` con el historial reescrito (`git filter-branch --subdirectory-filter`), así que sus commits no coinciden con los del repo local. Reglas:
- **Nunca** subir a GitHub el repo local completo.
- Si se trabaja en la Mac sobre el clon de GitHub, los cambios viven solo ahí. Hay que avisarlo para no perderlos al volver a la PC.

## 3. Fuentes de verdad, en orden

1. El caso (`2907_10433.pdf`): consignas a–e.
2. El Trabajo Escrito (TE), congelado y ya entregado a los sinodales.
3. `PLAN.md`.

Toda desviación del TE se valida con el usuario y se registra (ver `PLAN.md` §1.3 y los ADR).

## 4. Decisiones del usuario que no hay que volver a preguntar

- **Demo en la Mac y sin despliegue en la nube.** Vercel se descartó: contradice el TE §6.4.1 y ADR-005 (solo loopback), abre preguntas de seguridad y privacidad, y el respaldo ya existe sin internet (binario + `--selftest`).
- **Sin CI ni Wine** (desviación aprobada del TE §6.3.9). Se sustituye por la ejecución nativa en la Mac más la Capa 3 (`--selftest`).
- **Los umbrales del TE se miden y no bloquean** durante el sprint. La cobertura se informa con `make cover`.
- **No hay `docs/PLAN_PRUEBAS.md`:** el plan de pruebas son el TE §6.3 y el Anexo II.
- **Analizador léxico:**
  - stemmer ligero propio para agrupar variantes;
  - NFC escrita a mano y restringida al alfabeto (sin `golang.org/x/text`);
  - los números se reconocen pero no se catalogan.
  - **Propuesta de mejora** (no implementada; va en la réplica y en el TE §7.4): folios clave-valor como `TK-35699A`, que requerirían `-` en el alfabeto y un estado ALFANUM.
- **Corpus:**
  - Se corrigió la coherencia: un problema por comentario, montos coherentes y saludos y despedidas por categoría.
  - Después se corrigió la **separabilidad trivial**. Con el corpus `b810bad2…` el modelo sacaba 1.0, lo que violaba el TE §6.1.4. La solución fue un criterio registrado **antes** de cambiar plantillas (K1 Jaccard con stoplist, K2 colisión 5–10 %, K3 sonda NB en [0.80, 0.99]) y una **ambigüedad controlada** (10 oraciones compartidas entre los pares que el TE anticipa, más 30 plantillas de contexto cruzado inspiradas en `docs/data/grok_examples.txt`). No es ruido de etiquetas. Está documentado en **ADR-004 §6**.
  - Corpus actual: **`37f285ff…`**.
- **Revisiones manuales del usuario:**
  - L-1 se etiquetó por **16 contextos** (`docs/reports/l1_contextos.csv`). Hoy da 86.0 %.
  - Muestra de 250 y plantillas nuevas: el usuario las estaba revisando el 12-oct (`docs/reports/muestra_250_borrador.md` y `plantillas_nuevas.md`). Hay que registrar su constancia en `PLAN.md` §5.2 cuando la entregue.
- **UI 1.0.3** (la 1.0.3 agrega `/casos-prueba` con la matriz del Anexo II):
  - navegación en la misma pestaña, con "← Volver al clasificador";
  - `/catalogo` con filtros y el panel "Ver datos en JSON";
  - `/api` con la documentación de los 5 endpoints y un bloque por código de error (200, 400, 413, 405, 404 y 500, este con la explicación de que no ocurre en operación normal);
  - el pie no muestra la IP;
  - **sin Swagger**, porque rompería la regla de no usar recursos externos.

## 5. Cifras para la presentación (sección "Resultados")

| Medición | Valor |
|---|---|
| Corpus | 73,000 comentarios, semilla 10433, 30/20/20/20/10 |
| Catálogo | 6,196 lexemas; I-1 = 0 fallos, I-2 = 0 filas, I-3 = SHA idéntico |
| Modelo ganador | Regresión Logística Softmax (línea base, C = 1); CV 0.978 contra NBM 0.972 |
| Macro-F1 test temporal (CPM01) | **0.9795** (umbral 0.85) |
| Recall de Reconocimientos (CPM02) | **0.9819** (umbral 0.70) |
| Sin duplicados / sin ambiguos | 0.9826 / 0.9997 |
| **Plantillas no vistas** (cifra realista) | **0.835** |
| Paridad Python↔Go | 1,000/1,000 vectores idénticos bit a bit |
| L-1 / L-2 | 86.0 % / 91.4 % |
| `--selftest` / P99 | 0.05 s / ~0.1–0.2 ms (Linux) |

## 6. Respuestas preparadas para la réplica

- **"¿Por qué 0.98?"** Es el test temporal; el corpus es sintético y basado en plantillas. Con frases nunca vistas baja a 0.835, que es la estimación realista. Antes de la corrección daba 1.0 (trivialmente separable). Se diagnosticó, se fijó un criterio antes de cambiar los datos y se entrenó una sola vez.
- **"¿La ambigüedad no es ruido de etiquetas?"** No. La etiqueta nunca contradice el texto, es determinista y la tasa está declarada (8.95 %). El TE la anticipa en las causas primarias, en §5.2.2 y en el guion ("reconocimiento con tono ambiguo").
- **Cita incorrecta en el TE:** Hopcroft §3.4 es "Álgebra de las expresiones regulares". La complejidad O(n·m) está en **§4.3.2** (pertenencia) y en **Aho et al. §3.7.2–3.7.3** (algoritmo 3.22, fig. 3.37). El Teorema 2.11 prueba la equivalencia de la construcción de subconjuntos; el caso exponencial es §2.3.6. Respuesta acordada: *"tomé el número de sección de mis notas y, por la carga de trabajo durante el plazo de entrega, no lo verifiqué contra el libro; lo detecté al implementar el motor y lo corregí en ADR-002."* Dilo una vez y sin justificarte de más.
- **Simulación del AFN-ε:** `internal/regex/nfa.go`, función `simulate`, corresponde línea por línea a la fig. 3.37 de Aho et al.
- **Sin CI ni Wine:** la ejecución nativa en la Mac es evidencia más fuerte que un emulador, y la Capa 3 (`--selftest`) se conserva.
- **500 en `/api`:** las bases van embebidas y son de solo lectura. Si estuvieran dañadas, el servicio no arranca y `--selftest` lo detecta. El 500 solo cubre una falla extrema en ejecución.
- **"bajo" y "sobre" no aparecen en el corpus.** L-1 se midió sobre lo que el corpus contiene. Es un hallazgo de validez externa (TE §7.3).
- **El Anexo II del TE se ve mal en el PDF:** la última columna ("Salida esperada") se corta en el margen, y las barras `|` de las expresiones regulares partieron las celdas (CP-U01 era `"factur(a|as|…)"` y CP-U02 era `"factur(a|)"`). Es un problema de formato, no de contenido. No lo menciones tú. Si alguien lo nota: *"El Anexo II es una selección representativa; los criterios completos están en §6.3 y la matriz completa con su resultado está en el repositorio"*, y abres `/casos-prueba`. Dato a favor: es justo el problema que resuelve el escape `\|` de la gramática del TE §6.2.3.
- **`facturacion` sin tilde es un lexema aparte:** no hay plegado de acentos (ADR-003) y el stemmer no reconoce "-acion". Es una limitación documentada.

## 7. Preferencias de trabajo del usuario

- Comunicación en español y en lenguaje simple; explicar con ejemplos cuando algo es técnico.
- Planificar antes de implementar. Usar **subagentes Sonnet** para implementar y **Opus** para análisis delicados.
- Cada tarea en su rama, con commits rojo → verde. **Merge y release solo con su aprobación.**
- Revisar uno mismo los artefactos (muestras, capturas de la UI) **antes** de pasárselos al usuario.
- No usar herramientas de SQLite para VS Code; el usuario usa **DBeaver CE**.

## 8. Pendientes al 12-oct

1. Constancia del usuario: `plantillas_nuevas.md` (47 casillas) y `muestra_250_borrador.md`. Registrarla en `PLAN.md` §5.2.
2. ~~Confirmar en Safari~~: hecho; el usuario confirmó que la UI funciona en Safari en la Mac.
3. I-4 y la matriz de casos de prueba están hechas en la rama `fase5/casos-prueba` (pendiente de aprobación y merge): `make i4` (R² 0.9938, T(73k) 0.37 s en Linux; subcorpus = primeros N comentarios de `data/corpus.db`, sin `gencorpus -total` porque la batería K1–K3 está registrada para el corpus completo), `make casos` y la página `/casos-prueba` (29 casos del Anexo II + CPU14 como extensión).
4. Presentación (PPT) de 20 min: portada, diagnóstico, resolución de consignas y resultados (guía del Acuerdo 286), siguiendo el orden del TE. Usa las cifras de §5 y las respuestas de §6.
