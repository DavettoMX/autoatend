"""Builds docs/reports/evaluacion_modelo.md (Spanish). Deterministic: no timestamps, no timings."""

from __future__ import annotations

from .data import CATEGORIES, RECONOCIMIENTOS
from .models import GRIDS, LR_MAX_ITER, LR_TOL, N_FOLDS, SEED

PARITY_MARKER = "<!-- PARIDAD -->"

FAM = {"nbm": "Naive Bayes multinomial", "rls": "Regresión logística softmax"}
STRAT = {
    "pesos_balanceados": "pesos balanceados",
    "sobremuestreo": "sobremuestreo de Reconocimientos",
    "linea_base": "línea base",
}
THRESH_F1 = 0.85
THRESH_RECALL = 0.70


def hp(params: dict) -> str:
    return ", ".join(f"{k} = {v:g}" for k, v in params.items())


def f(x: float, nd: int = 4) -> str:
    return f"{x:.{nd}f}"


def verdict(value: float, thr: float) -> str:
    return "CUMPLE" if value >= thr else "NO CUMPLE"


def class_table(ev: dict) -> str:
    rows = ["| Categoría | Precisión | Recall | F1 | Soporte |", "|---|---|---|---|---|"]
    for c in ev["por_clase"]:
        rows.append(f"| {c['nombre']} | {f(c['precision'])} | {f(c['recall'])} | {f(c['f1'])} | {c['soporte']:,} |")
    return "\n".join(rows)


def confusion(ev: dict) -> str:
    head = "| Real \\ Predicho | " + " | ".join(c[1] for c in CATEGORIES) + " |"
    rows = [head, "|---|" + "---|" * len(CATEGORIES)]
    for i, (_, name) in enumerate(CATEGORIES):
        rows.append(f"| {name} | " + " | ".join(f(v) for v in ev["confusion_normalizada"][i]) + " |")
    return "\n".join(rows)


def confusion_counts(ev: dict) -> str:
    head = "| Real \\ Predicho | " + " | ".join(c[1] for c in CATEGORIES) + " |"
    rows = [head, "|---|" + "---|" * len(CATEGORIES)]
    for i, (_, name) in enumerate(CATEGORIES):
        rows.append(f"| {name} | " + " | ".join(f"{v:,}" for v in ev["confusion_conteos"][i]) + " |")
    return "\n".join(rows)


def metrics_block(title: str, ev: dict) -> str:
    return (
        f"#### {title}\n\n"
        f"n = {ev['n']:,}; exactitud = {f(ev['accuracy'])}; **Macro-F1 = {f(ev['macro_f1'])}**; "
        f"**recall de Reconocimientos = {f(ev['recall_reconocimientos'])}**.\n\n"
        f"{class_table(ev)}\n\n"
        f"Matriz de confusión normalizada por fila (cada fila suma 1):\n\n{confusion(ev)}\n\n"
        f"Conteos:\n\n{confusion_counts(ev)}\n"
    )


def build_report(c: dict) -> str:
    vocab = c["vocab"]
    w, alt = c["winner"], c["alt"]
    fin, fa = c["finals"]["principal"], c["finals"]["alternativo"]
    ev, evn = fin["eval"], fin["eval_nd"]
    n_train = vocab.n_train
    L: list[str] = []
    add = L.append

    add("# Evaluación del modelo (Bloque C, consigna e)\n")
    add(
        "Reporte generado por `make train` (`ml/autoatend_ml`). Es determinista: dos corridas producen el mismo archivo. "
        f"Las cifras salen de `data/corpus.db` (SHA `{c['corpus_sha'][:8]}…`), `data/df_train.json` y `data/corpus_ambiguos.json`; ninguna se ajustó sobre test. "
        "El generador de este corpus cumple el criterio de separabilidad K1–K3, registrado en `PLAN.md` §3 antes de cambiar las plantillas; "
        "el modelo se entrenó una sola vez, después de que la batería pasó.\n"
    )

    add("## 1. Datos y vocabulario\n")
    add(f"- Corte temporal (`db.SplitCutoff`, misma consulta SQL): `{vocab.cutoff}`. Train = `fecha < corte`: **{n_train:,}** comentarios; test = `fecha >= corte`: **{c['n_test']:,}**.")
    add(f"- Catálogo de `resultados.db` (corpus completo): {c['n_catalogue']:,} lexemas. Vocabulario del modelo (df > 0 en train): **{len(vocab.lexemes):,}** lexemas.")
    add(f"- Lexemas con **df = 0 en train, excluidos del modelo: {c['df_cero']:,}** (aparecen solo en test; el TE §6.4.2 calcula el vocabulario y el IDF con train).")
    add(f"- Verificación cruzada: la frecuencia de documento recalculada en Python con la salida de `cmd/tokenizar` coincide con `df_train.json` ({c['df_check'][0]} discrepancias, {c['df_check'][1]} lexemas de más).")
    add(f"- Comentarios de test con vector vacío (ningún lexema del vocabulario): {c['n_empty_test']:,}.")
    add(f"- Duplicados exactos: **{c['n_dup']:,} de {c['n_test']:,}** comentarios de test ({100 * c['n_dup'] / c['n_test']:.1f} %) tienen un texto idéntico en train; el subconjunto sin duplicados tiene **{c['n_nodup']:,}**.")
    add(f"- Subconjunto ambiguo (K2): {c['n_amb_total']:,} comentarios del corpus cuya oración base aparece, como bolsa de tokens, en dos categorías (plantillas compartidas: ambigüedad controlada, sin ruido de etiquetas). En test son **{c['n_amb_test']:,}**; el test sin ellos tiene **{c['n_noamb']:,}**.\n")

    add("## 2. Método\n")
    add("- **Tokenización**: `cmd/tokenizar` (Go) invocado una vez como subproceso (TE §6.3.6). Cada forma → `lexer.Stem` → lexema canónico del catálogo; números y stoplist fuera; raíces fuera del catálogo se descartan.")
    add("- **Características** (TE §5.1.3, §6.4.2): TF = conteo crudo; IDF suavizado `idf = ln((1 + N) / (1 + df)) + 1` con N y df de train (`df_train.json`); vector normalizado en L2. Los productos y la suma de cuadrados se hacen en Python puro, en orden ascendente de índice, para que Go reproduzca los bits.")
    alphas = ", ".join(format(p["alpha"], "g") for p in GRIDS["nbm"])
    cs = ", ".join(format(p["C"], "g") for p in GRIDS["rls"])
    add(f"- **Familias**: Naive Bayes multinomial (suavizado de Laplace/Lidstone, α ∈ {{{alphas}}}) y regresión logística softmax multinomial con penalización L2 (C ∈ {{{cs}}}). La regresión usa `solver='saga'` (gradiente promedio estocástico, de la familia del descenso de gradiente estocástico), `max_iter = {LR_MAX_ITER}`, `tol = {LR_TOL:g}`.")
    add("- **Desbalance (30/20/20/20/10 %)**, tres estrategias: (a) pesos balanceados (`class_weight='balanced'`; en Naive Bayes, `sample_weight` de `compute_sample_weight('balanced')`); (b) sobremuestreo aleatorio de Reconocimientos hasta el tamaño mediano de clase (20 %), con semilla fija, aplicado solo a las particiones de entrenamiento de cada pliegue y al ajuste final, nunca a validación ni a test; (c) línea base sin ajuste.")
    add(f"- **Selección**: validación cruzada estratificada de {N_FOLDS} particiones (`StratifiedKFold(5, shuffle=True, random_state={SEED})`) sobre train. Gana la celda con mayor Macro-F1 promedio; en un empate, la de menor desviación y luego la primera de la rejilla. `n_jobs = 1`, semillas fijas.")
    add("- **Limitación declarada**: el vocabulario y el IDF se calculan una vez con todo train (como pide el TE), no por pliegue. Es una fuga no supervisada de estadísticas de documento hacia los pliegues de validación; puede inflar ligeramente el Macro-F1 de CV y no afecta a test.")
    add("- **Modelo lineal único** para las dos familias: `z_k = b_k + Σ_i w_{k,i}·x_i`, softmax con resta del máximo. Naive Bayes: `w = feature_log_prob_`, `b = class_log_prior_`; regresión: `coef_`, `intercept_`.\n")

    add("## 3. Validación cruzada (train)\n")
    add("Mejor hiperparámetro de cada celda; Macro-F1 promedio ± desviación estándar muestral de las 5 particiones.\n")
    add("| Familia | Estrategia | Mejores hiperparámetros | Macro-F1 CV (media ± desv.) |")
    add("|---|---|---|---|")
    for cell in c["cells"]:
        mark = " **(ganador)**" if cell is w else (" (mejor de la otra familia)" if cell is alt else "")
        add(f"| {FAM[cell['familia']]} | {STRAT[cell['estrategia']]} | {hp(cell['hiperparametros'])} | {f(cell['f1_media'], 5)} ± {f(cell['f1_desv'], 5)}{mark} |")
    add("")
    add("Rejilla completa:\n")
    add("| Familia | Estrategia | Hiperparámetros | Macro-F1 CV (media ± desv.) | Convergió |")
    add("|---|---|---|---|---|")
    for g in c["grid"]:
        conv = "n/a" if g["familia"] == "nbm" else ("sí" if g["convergio"] else "no")
        add(f"| {g['familia']} | {g['estrategia']} | {hp(g['hiperparametros'])} | {f(g['f1_media'], 5)} ± {f(g['f1_desv'], 5)} | {conv} |")
    add("")
    near = [x for x in c["cells"] if x is not w and w["f1_media"] - x["f1_media"] <= w["f1_desv"]]
    allf = [g["f1_media"] for g in c["grid"]]
    if max(allf) - min(allf) < 0.001:
        add(
            f"Las {len(c['grid'])} combinaciones de la rejilla quedan entre {f(min(allf), 5)} y {f(max(allf), 5)} de Macro-F1 promedio: "
            "el corpus es tan separable que la estrategia de desbalance y el hiperparámetro cambian el resultado en la cuarta o quinta cifra decimal.\n"
        )
    else:
        add(f"Las {len(c['grid'])} combinaciones de la rejilla quedan entre {f(min(allf), 5)} y {f(max(allf), 5)} de Macro-F1 promedio.\n")
    if near:
        add(
            f"**Empate o diferencia menor a una desviación estándar** entre el ganador y {len(near)} celda(s) más: la elección entre ellas "
            "no es estadísticamente firme y la decide el criterio de desempate declarado en el método (menor desviación, luego primera de la rejilla).\n"
        )
    edge = None
    if w["familia"] == "rls" and w["hiperparametros"]["C"] == max(p["C"] for p in GRIDS["rls"]):
        edge = "C = 10 es el extremo superior de la rejilla"
    if w["familia"] == "nbm" and w["hiperparametros"]["alpha"] in (
        min(p["alpha"] for p in GRIDS["nbm"]),
        max(p["alpha"] for p in GRIDS["nbm"]),
    ):
        edge = "α está en un extremo de la rejilla"
    if edge:
        add(f"El ganador queda en el borde de la rejilla ({edge}); la rejilla es la fijada por el diseño y no se amplió.\n")

    add("## 4. Ganador\n")
    add(f"**{FAM[w['familia']]}**, estrategia **{STRAT[w['estrategia']]}**, {hp(w['hiperparametros'])}; Macro-F1 de CV {f(w['f1_media'], 5)} ± {f(w['f1_desv'], 5)}. Reajustado con todo train y exportado como `data/modelo.json`.")
    add(f"El mejor de la otra familia ({FAM[alt['familia']]}, {STRAT[alt['estrategia']]}, {hp(alt['hiperparametros'])}; CV {f(alt['f1_media'], 5)} ± {f(alt['f1_desv'], 5)}) se exporta como `data/modelo_alternativo.json`.\n")

    add("## 5. Evaluación final en test (una sola vez, sin ajustar nada)\n")
    eva = fin["eval_na"]
    add("| Criterio | Umbral | Test completo | Test sin duplicados de train | Test sin subconjunto ambiguo |")
    add("|---|---|---|---|---|")
    add(f"| CPM01 Macro-F1 | ≥ {THRESH_F1} | {f(ev['macro_f1'])} ({verdict(ev['macro_f1'], THRESH_F1)}) | {f(evn['macro_f1'])} ({verdict(evn['macro_f1'], THRESH_F1)}) | {f(eva['macro_f1'])} ({verdict(eva['macro_f1'], THRESH_F1)}) |")
    add(f"| CPM02 recall de Reconocimientos | ≥ {THRESH_RECALL} | {f(ev['recall_reconocimientos'])} ({verdict(ev['recall_reconocimientos'], THRESH_RECALL)}) | {f(evn['recall_reconocimientos'])} ({verdict(evn['recall_reconocimientos'], THRESH_RECALL)}) | {f(eva['recall_reconocimientos'])} ({verdict(eva['recall_reconocimientos'], THRESH_RECALL)}) |")
    add(f"| n | | {ev['n']:,} | {evn['n']:,} | {eva['n']:,} |")
    add("")
    add(metrics_block(f"Modelo ganador, test completo ({FAM[w['familia']]})", ev))
    add(metrics_block("Modelo ganador, test sin duplicados exactos de train", evn))
    add(metrics_block("Modelo ganador, test sin el subconjunto ambiguo", eva))
    add("### Modelo alternativo (mejor de la otra familia)\n")
    add("| Modelo | Test: Macro-F1 | Test: recall Reconocimientos | Sin duplicados: Macro-F1 | Sin duplicados: recall Reconocimientos | Sin ambiguos: Macro-F1 | Sin ambiguos: recall Reconocimientos |")
    add("|---|---|---|---|---|---|---|")
    for tag, cell, fx in (("modelo.json", w, fin), ("modelo_alternativo.json", alt, fa)):
        e1, e2, e3 = fx["eval"], fx["eval_nd"], fx["eval_na"]
        add(f"| `{tag}` ({cell['familia']}, {cell['estrategia']}) | {f(e1['macro_f1'])} | {f(e1['recall_reconocimientos'])} | {f(e2['macro_f1'])} | {f(e2['recall_reconocimientos'])} | {f(e3['macro_f1'])} | {f(e3['recall_reconocimientos'])} |")
    add("")

    nb = [(k, x) for k, x in (("principal", fin), ("alternativo", fa)) if x["model"]["familia"] == "nbm"]
    add("### Nota de cordura (corpus sintético)\n")
    if nb:
        k, x = nb[0]
        v = x["eval"]["macro_f1"]
        if v > 0.99:
            add(f"**El Naive Bayes multinomial alcanza Macro-F1 = {f(v)} (> 0.99) en test.** Es una señal de que el corpus sintético, generado a partir de plantillas por categoría, es fácil de separar. Los umbrales del TE se cumplen con holgura, pero esa holgura describe al corpus y no la dificultad del problema real de TelArg.\n")
        else:
            add(f"El Naive Bayes multinomial alcanza Macro-F1 = {f(v)} en test, por debajo de 0.99.\n")
    else:
        add("No hay modelo Naive Bayes entre los exportados.\n")
    add(f"Además, el {100 * c['n_dup'] / c['n_test']:.1f} % del test repite textos de train; por eso las métricas se reportan también sin ellos (sección 5).\n")
    add("**F1 con plantillas no vistas (solo reporte, no es compuerta ni umbral).** El corte temporal no excluye plantillas: todas aparecen en train y en test. "
        "Para medir cuánto depende el modelo de haber visto cada plantilla, `docs/reports/plantillas_no_vistas.md` reentrena con el 80 % de las plantillas "
        "de cada categoría y evalúa con el 20 % restante, en 5 pliegues. Es un análisis aparte, que no participa en la selección ni en la evaluación de esta sección.\n")

    add("## 6. Equivalencia del escalador lineal con scikit-learn\n")
    add("| Modelo | Máx. \\|p_lineal − predict_proba\\| sobre todo test | Acuerdo de clase |")
    add("|---|---|---|")
    for tag, fx in (("modelo.json", fin), ("modelo_alternativo.json", fa)):
        add(f"| `{tag}` ({fx['model']['familia']}) | {fx['linealidad']['max_dif_prob']:.2e} | {fx['linealidad']['acuerdo_clase']:.6f} |")
    add("\n`softmax(b + W·x)` con los pesos exportados reproduce `predict_proba` de scikit-learn (Naive Bayes: `softmax` de la verosimilitud conjunta = posterior).\n")

    add("## 7. Paridad Python ↔ Go (CPA06)\n")
    add(PARITY_MARKER + "\n")

    add("## 8. Comentarios de la demostración\n")
    add("`data/gold5.json`: un comentario de test por categoría, sin duplicado en train, fuera del subconjunto ambiguo (nunca de una plantilla compartida), longitud legible, clasificado correctamente con confianza ≥ 0.80 (elegido al azar con la semilla 10433 entre los elegibles).\n")
    add("| Categoría | id | Confianza | Texto |")
    add("|---|---|---|---|")
    for g in c["gold"]:
        add(f"| {g['categoria']} | {g['id']} | {f(g['confianza_python'], 3)} | {g['texto']} |")
    add("")
    add("Candidatos para la demostración oral (no están en `gold5.json`):\n")
    for key, label in (("claro_falta_de_servicio", "Caso claro, Falta de Servicio"), ("frontera_reconocimiento", "Caso frontera, Reconocimientos")):
        d = c["demo"][key]
        add(f"- **{label}** (id {d['id']}): «{d['texto']}». Real: {d['categoria_real']}; predicción: {d['prediccion']} con {f(d['confianza'], 3)}; segunda: {d['segunda']} con {f(d['confianza_segunda'], 3)}. Lexemas: {', '.join(d['lexemas'])}.")
    add("")
    add("El caso frontera es un Reconocimiento del subconjunto ambiguo (una plantilla compartida con Quejas, el «reconocimiento con tono ambiguo» del guion del TE §6.4.3): entre los del test, el de menor margen entre las dos clases más probables, acertado si alguno lo está. El caso claro sale del test sin duplicados ni ambiguos.\n")

    add("## 9. Reproducción\n")
    add("```bash\nmake train        # compila cmd/tokenizar, entrena, escribe modelo.json, modelo_alternativo.json, gold5.json, testdata/paridad_1000.json y este reporte\nmake train-sha    # dos entrenamientos independientes producen el mismo SHA-256\n```\n")
    return "\n".join(L)
