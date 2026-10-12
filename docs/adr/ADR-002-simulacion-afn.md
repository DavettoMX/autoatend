# ADR-002: simulación directa del AFN-ε por conjuntos de estados

**Estado:** aceptada (Fase 2, Bloque B). Es la decisión que el Trabajo Escrito (TE) §6.2.3 promete documentar.
**Implementación:** `internal/regex` (`ast.go`, `parse.go`, `nfa.go`) y la validación de ER de `internal/sweep` (`persist.go`).
**Pruebas:** CPU01 a CPU08, CPP01 (I-1 sobre `data/resultados.db`) e integración 4b.

## Contexto

La consigna b pide un analizador léxico cuyas expresiones regulares (ER) se interpreten con un motor propio. El TE §6.2.3 fija la gramática y §6.2.4 la construcción de Thompson, y deja abierto cómo se ejecuta el autómata resultante: se puede (a) simular el AFN-ε directamente o (b) convertirlo a un AFD con la construcción de subconjuntos y recorrer el AFD. El motor corre en dos lugares: el round-trip de cada ER antes de persistirla (6,720 lexemas en el catálogo) y la verificación de I-1 sobre todas las variantes observadas (7,439).

Restricciones: Go puro sin CGo, prohibido `import "regexp"` y cualquier generador de parsers (`make lint-regexp`).

## Decisión

El AFN-ε de Thompson se **simula directamente por conjuntos de estados**, sin construir un AFD (Thompson 1968; Aho et al. 2006, §3.7.2 y algoritmo 3.22; Hopcroft et al. 2007, §2.5.3 para la clausura respecto de ε y §4.3.2 para la comprobación de pertenencia):

1. El estado inicial del recorrido es `S₀ = ε-cerradura({inicio})`.
2. Por cada runa `c` de la cadena, `S' = ε-cerradura(movida(S, c))`.
3. La cadena se acepta si y solo si el conjunto final contiene el estado de aceptación. La coincidencia es **de cadena completa**: una ER acepta una variante si y solo si toda la variante pertenece a L(ER).

La ε-cerradura se calcula con una pila explícita y una marca por estado (sellada con un número de generación), así que un ciclo ε como el de `(a*)*` no puede ciclar ni desbordar la pila, y cada estado entra a lo más una vez por paso. No hay retroceso.

**Costo.** Para una cadena de `n` runas y un AFN de `m` estados: `O(n·m)` en tiempo (cada paso examina a lo más las transiciones de los `m` estados) y `O(m)` en espacio (el conjunto, la pila y las marcas). La construcción de Thompson es lineal: cada operador de la ER añade a lo más 2 estados, de modo que `m ≤ 2·|r|` (`TestCPU05_ThompsonStateBound` lo comprueba, y el catálogo completo en la integración 4b).

**Parámetros del motor** (registrados aquí porque el TE no los fija):

- Profundidad máxima de anidamiento de paréntesis, `regex.MaxDepth = 200`. El analizador sintáctico y la construcción recursivos gastan unos pocos marcos por nivel, por lo que el límite acota la pila para cualquier entrada (una cadena de 2 millones de `(` se rechaza en el nivel 201). El caso de esfuerzo del TE, 50 niveles en 500 caracteres, queda muy por debajo.
- `LETRA` es cualquier letra Unicode y `DIGITO` un dígito decimal; el alfabeto del lexer es `a-z` y `áéíóúüñ`.
- Los metacaracteres que se escapan son `| * + ? ( ) [ ] \ -`, el mismo conjunto que escapa `lexer.GenerateER`; un `-` sin escapar fuera de una clase es un error.
- Una alternancia de `k` ramas se construye plegando por pares, de izquierda a derecha, con la regla binaria del TE (un inicio y un aceptador nuevos por par), así que `a|b|c` tiene 10 estados.
- `+` y `?` se derivan de la regla de `*`: `+` omite la ε del inicio nuevo al aceptador nuevo y `?` omite la ε del aceptador de N de regreso a su inicio.
- Una clase de caracteres es un fragmento de dos estados con una transición etiquetada por rango.

## Alternativas descartadas

| Alternativa | Por qué se descarta |
|---|---|
| **AFN → AFD por construcción de subconjuntos**, y recorrer el AFD (Aho et al. 2006, §3.7.1; Hopcroft et al. 2007, §2.3.5, Teorema 2.11; el caso exponencial está en §2.3.6, "Un caso desfavorable para la construcción de subconjuntos", y el costo de las conversiones en §4.3.1) | El AFD hace cada comparación en `O(n)`, pero la construcción cuesta hasta `O(2^m)` en el peor caso y debe repetirse para cada ER. Aquí las ER son cortas y numerosas (una por lexema) y cada una se ejecuta contra unas pocas variantes (7,439 comparaciones para 6,720 ER, es decir 1.1 por ER), de modo que el costo de construcción no se amortiza. Además añade código (estados del AFD como conjuntos canónicos, memoización) que el TE no pide y que habría que defender y probar. |
| **AFD perezoso** (construir los subconjuntos al vuelo y guardarlos en caché) | Mantiene el peor caso lineal, pero complica el motor con una caché por ER y exige sincronización si el NFA se comparte entre goroutines. Con cadenas de pocas runas la caché casi nunca se reutiliza. |
| **Retroceso** (como el motor de Perl o de PCRE) | Es exponencial en el peor caso: `(a?){n}a{n}` contra `aⁿ`. `TestCPU08_NoBacktrackingBlowup` fija n = 60 y falla si tarda más de un segundo. No ofrece garantía de `O(n·m)`, que es la que promete el TE. |
| **Usar `regexp` de la biblioteca estándar** | Prohibido por la consigna b (el analizador léxico debe ser propio) y por `make lint-regexp`. |
| **Derivadas de Brzozowski** | Es una técnica válida, pero no es la construcción de Thompson que fija el TE §6.2.4, y el TE exige un AFN-ε con un estado inicial y uno de aceptación por fragmento. |

## Verificación

- **Corrección:** `TestCPU04` a `TestCPU07` y `TestMatchSemantics` (cadena completa, cerraduras, clases, escapes, ciclos ε). `TestCPP01_GeneratorRoundTripRandom` genera 2,000 conjuntos de variantes aleatorios sobre un alfabeto que mezcla letras y todos los metacaracteres, y comprueba que la ER de `GenerateER` acepta exactamente las variantes.
- **Complejidad:** `TestCPU08_LinearCost` cuenta de forma determinista los estados visitados y las transiciones examinadas, y exige que el trabajo crezca 10 veces (con tolerancia de 8 a 12) entre 10, 100, 1,000 y 10,000 caracteres y que no pase de `8·n·m`. `TestCPU08_LinearTime` mide el tiempo real con la misma escala (se omite con `-short`) y falla si algún paso crece más de 40 veces, ya que un algoritmo cuadrático crecería 100.
- **Invariante I-1:** `TestCPP01_I1_RoundTripCatalogue` abre `data/resultados.db` en solo lectura, compila la ER de cada lexema y comprueba que acepta todas sus variantes: 6,720 lexemas, 7,439 variantes, 0 fallos. El barrido hace el mismo round-trip antes de persistir cada ER (TE Etapa 6): si una ER no compila o rechaza una variante, el lexema cae a su forma canónica y `Result.ERFallbacks` lo cuenta (hoy 0, y por eso `resultados.db` no cambió de SHA).
- **Tamaño del AFN sobre el catálogo:** máximo 46 estados, media 14.4.

## Consecuencias

- El motor es pequeño (`parse.go`, `nfa.go`) y cada paso del TE tiene su prueba con el ID del Anexo II.
- Una ER compilada es inmutable: un mismo `*NFA` se puede usar desde varias goroutines (cada `Match` reserva su propia memoria de trabajo, `O(m)`).
- Cada llamada a `Match` reserva el conjunto y las marcas. Si en la Fase 4 la ruta de inferencia llegara a ejecutar ER en caliente, el remedio es reutilizar esa memoria con un `sync.Pool`; hoy no hace falta, porque la clasificación no ejecuta ER (solo el barrido y la verificación lo hacen).
- Si en el futuro hubiera pocas ER y millones de cadenas por ER, convendría la conversión a AFD (u otra técnica); no es el caso de este sistema.

## Nota sobre las referencias del Trabajo Escrito

El TE (§5.1.1, §6.2.4 y la tabla de §6.5) atribuye la complejidad `O(n·m)` de la simulación del AFN-ε a Hopcroft et al. (2007, §3.4), y la explosión exponencial de la conversión a AFD al Teorema 2.11. Cotejado con la edición en español del libro:

| Cita del TE | Contenido real | Referencia correcta |
|---|---|---|
| §3.4 | "Álgebra de las expresiones regulares" (leyes de equivalencia) | §4.3.2, "Cómo comprobar la pertenencia a un lenguaje regular" (p. 129), y Aho et al. (2006), §3.7.2–§3.7.3 |
| Teorema 2.11 | Corrección de la construcción de subconjuntos: L(D) = L(N) (§2.3.5) | Correcto para la equivalencia; el peor caso exponencial es §2.3.6 (p. 54), y el costo de convertir AFN en AFD, §4.3.1 (p. 126) |

El TE está congelado y no se modifica. Este ADR usa las referencias correctas. El contenido técnico del TE (algoritmo, complejidad y argumento de diseño) no cambia: solo cambia el número de sección que lo respalda.

## Referencias

- Thompson, K. (1968). Regular expression search algorithm. *Communications of the ACM*, 11(6), 419-422.
- Aho, A. V., Lam, M. S., Sethi, R. y Ullman, J. D. (2006). *Compilers: Principles, Techniques, and Tools* (2.ª ed.), §3.7.
- Hopcroft, J. E., Motwani, R. y Ullman, J. D. (2007). *Introduction to Automata Theory, Languages, and Computation* (3.ª ed.; ed. en español: *Introducción a la teoría de autómatas, lenguajes y computación*, Pearson), §2.3.5–§2.3.6, §2.5.3 y §4.3.1–§4.3.2.
