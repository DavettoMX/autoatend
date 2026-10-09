# CENEVAL Case 10433 - Contexto Arquitectónico (AutoAtend)

Este proyecto está bajo un escrutinio estricto para una evaluación del CENEVAL. El entorno objetivo de despliegue ("Deployment Envelope") es Windows x86-64 cerrado, sin instalación previa y sin conexión a internet.

**REGLA CRÍTICA PARA EL AGENTE:**
Antes de comenzar cualquier tarea de código, estructuración de carpetas o tomar decisiones de diseño en una nueva sesión, **DEBES** leer el índice de arquitectura que reside en Obsidian:
`/home/davetto/Documentos/Obsidian Vault/Proyecto Acreditacion Ceneval/Antigravity/00_Indice_Arquitectura.md`

**Restricciones inquebrantables de diseño rápido:**
1. **Runtime:** Estrictamente **Golang**. No proponer Python, Node.js ni Java para el runtime de producción.
2. **Base de Datos:** SQLite embebido (usar `modernc.org/sqlite`).
3. **Analizador Léxico:** Debe programarse a mano en Go (cero dependencias de NLP).
4. **Clasificador IA:** Entrenado offline en Python, pero los pesos deben exportarse como JSON puro ("Plain Data") e inferirse matemáticamente en Go. No usar ONNX.
