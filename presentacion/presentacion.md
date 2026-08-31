---
marp: true
paginate: true
size: 16:9
style: |
  section {
    font-family: 'Helvetica Neue', Arial, sans-serif;
    background: #ffffff;
    color: #1e293b;
    padding: 50px 65px;
  }
  h1 {
    color: #4338ca;
    font-size: 1.9em;
    border-bottom: 3px solid #4338ca;
    padding-bottom: 0.2em;
  }
  h2 {
    color: #4338ca;
  }
  section.lead {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    background: linear-gradient(135deg, #312e81 0%, #4338ca 100%);
    color: white;
  }
  section.lead h1 {
    color: white;
    border: none;
    font-size: 2.6em;
  }
  section.lead p {
    color: #c7d2fe;
    font-size: 1.2em;
  }
  section.section {
    background: #4338ca;
    color: white;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  section.section h1 {
    color: white;
    border: none;
    font-size: 2.4em;
    text-align: center;
  }
  code {
    background: #f1f5f9;
    color: #4338ca;
  }
  pre {
    background: #1e293b;
  }
  table {
    font-size: 0.78em;
  }
  th {
    background: #4338ca;
    color: white;
  }
  .flow {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    margin-top: 35px;
    flex-wrap: nowrap;
  }
  .box {
    background: #eef2ff;
    border: 2px solid #4338ca;
    border-radius: 10px;
    padding: 10px 8px;
    text-align: center;
    font-size: 0.62em;
    font-weight: 600;
    color: #312e81;
    line-height: 1.35;
    flex: 1 1 0;
  }
  .arrow {
    font-size: 1.3em;
    color: #4338ca;
    font-weight: bold;
    flex: 0 0 auto;
  }
  .note {
    font-size: 0.7em;
    color: #64748b;
    margin-top: 25px;
  }
  .big {
    font-size: 1.4em;
    text-align: center;
    color: #4338ca;
    font-weight: 700;
    margin-top: 40px;
  }
  ul, ol {
    font-size: 0.95em;
  }
---

<!-- _class: lead -->

# Delivery Effort Estimator

Cómo funciona el motor de estimación de esfuerzo de delivery

<p>Framework AI-native · basado en evidencia · predicción, no compromiso</p>

---

## El problema

- Los equipos hoy combinan **trabajo humano y trabajo de agentes de IA** en cada feature.
- Story Points y horas fueron diseñados para equipos 100% humanos — **no distinguen** cuánto hace una persona y cuánto un agente.
- Las estimaciones se hacen "a ojo", sin conexión con lo que realmente pasó después.

<p class="big">¿Cómo predecir el esfuerzo de un feature de forma reproducible, auditable y medible contra la realidad?</p>

---

## Principios que guían el diseño

- **Evidencia, no intuición** — toda estimación se basa en la especificación/planning real, no en una corazonada.
- **Effort ≠ Duration ≠ Cost** — cuánto trabajo, cuánto tiempo y cuánto cuesta en IA son tres cosas distintas y se reportan por separado.
- **Es una predicción, no un compromiso** — se guarda para compararla después contra lo que realmente pasó.
- **El LLM no es el motor de estimación** — un LLM puede ayudar a describir el trabajo, pero el cálculo final siempre es una fórmula determinística y versionada.

---

<!-- _class: section -->

# El flujo completo

---

## De la especificación a la predicción

<div class="flow">
  <div class="box">Spec-Kit<br/>spec + plan + tasks</div>
  <div class="arrow">→</div>
  <div class="box">Plugin Claude Code<br/>deriva 7 features</div>
  <div class="arrow">→</div>
  <div class="box">Motor de Estimación<br/>determinístico</div>
  <div class="arrow">→</div>
  <div class="box">Predicción<br/>DEU · Riesgo · Costo</div>
</div>

<div class="flow" style="margin-top:45px;">
  <div class="box">Se hace<br/>el trabajo</div>
  <div class="arrow">→</div>
  <div class="box">Outcome Record<br/>lo que pasó</div>
  <div class="arrow">→</div>
  <div class="box">Error Report<br/>predicción vs. real</div>
  <div class="arrow">→</div>
  <div class="box">Calibración<br/>futura</div>
</div>

<p class="note">Nada se sobrescribe nunca: cada predicción y cada resultado quedan guardados tal cual se produjeron.</p>

---

## Dos decisiones de diseño clave

1. **Las features pueden venir de dos lugares** — una persona las escribe a mano en JSON, o el plugin de Claude Code las deriva automáticamente leyendo la spec de la feature. **El motor recibe exactamente lo mismo en ambos casos.**

2. **El motor nunca llama a un LLM.** Sin importar quién generó las 7 features, el cálculo que las convierte en una predicción es una fórmula fija — así la predicción es siempre reproducible y comparable en el tiempo.

---

<!-- _class: section -->

# Las 7 dimensiones de entrada

---

## Todo empieza con 7 números (0.0 a 1.0)

| Dimensión | Qué mide |
|---|---|
| `context_complexity` | Cuánto del sistema hay que entender para hacer el cambio |
| `domain_complexity` | Qué tan compleja es la lógica de negocio de fondo |
| `integration_complexity` | Cuántos otros sistemas toca (APIs, bases de datos, eventos) |
| `verification_complexity` | Qué tan difícil es demostrar que el resultado es correcto |
| `human_decision_load` | Cuántas decisiones reales quedan en manos de una persona |
| `ai_execution_complexity` | Cuánto razonamiento/iteración necesita el agente de IA |
| `uncertainty` | Cuánto sigue sin resolverse o sin saberse |

---

## ¿De dónde salen estos 7 números?

El plugin de Claude Code los deriva automáticamente leyendo `spec.md`, `plan.md` y `tasks.md` de la feature, siguiendo una **rúbrica documentada** (no intuición libre):

- Cada dimensión tiene anclas concretas (0.0 / 0.5 / 1.0) con señales específicas a buscar.
- Cada puntaje viene acompañado de una **justificación** que cita el contenido real de la spec.
- Si falta información, se puntúa de forma conservadora y se marca como supuesto — nunca se inventa.

<p class="note">Resultado: nada de "caja negra" — cualquiera puede auditar por qué se puntuó así.</p>

---

<!-- _class: section -->

# Cómo se calcula la predicción

---

## Paso 1 — cuatro esfuerzos independientes

Cada uno dominado por la dimensión que le corresponde (evidencia, no promedio parejo):

```text
Human Effort         = 0.45·human_decision_load + 0.30·context_complexity
                      + 0.15·domain_complexity  + 0.10·uncertainty

AI Effort             = 0.50·ai_execution_complexity + 0.25·context_complexity
                      + 0.15·domain_complexity  + 0.10·uncertainty

Verification Effort   = 0.55·verification_complexity + 0.20·integration_complexity
                      + 0.15·uncertainty        + 0.10·domain_complexity

Integration Effort    = 0.70·integration_complexity + 0.15·domain_complexity
                      + 0.15·uncertainty
```

---

## Paso 2 — Delivery Effort Unit (DEU)

```text
DEU = 10 × ( 0.35·HumanEffort + 0.20·AIEffort
           + 0.25·VerificationEffort + 0.20·IntegrationEffort )
```

- Escala relativa propia (0–10 aprox.), **no es horas ni Story Points**.
- Su significado real se aprende con el tiempo, comparando predicciones contra resultados reales.

---

## Paso 3 — Confianza (deliberadamente conservadora)

```text
Confidence = clamp(1 − uncertainty, 0.05, 0.70)
```

- Tope máximo de **70%** en esta primera versión del modelo (`v1-linear`).
- Aunque la incertidumbre reportada sea cero, el modelo **nunca** dice "estoy 100% seguro" — todavía no hay historial que lo respalde.

---

## Paso 4 — Rango de predicción (no un número falso-preciso)

```text
P50 = DEU
spread = 0.25 + 0.75 × (1 − Confidence)
P80 = DEU × (1 + spread)
```

<p class="big">"Lo más probable es X, pero razonablemente podría llegar a Y."</p>

---

## Pasos 5 y 6 — Riesgo, Duración y Costo

- **Riesgo** = directamente la `uncertainty` → etiqueta baja / media / alta.
- **Duración esperada** = `0.2 × DEU` días — una predicción **separada**, no el effort disfrazado de tiempo.
- **Costo de IA esperado** = `0.3 × AIEffort × DEU` USD — dimensión económica **independiente** del esfuerzo.

<p class="note">Effort, duración y costo nunca se mezclan en un solo número (principio de diseño).</p>

---

<!-- _class: section -->

# Ejemplo práctico

---

## Input: una feature real

```json
{
  "context_complexity": 0.72,
  "domain_complexity": 0.41,
  "integration_complexity": 0.83,
  "verification_complexity": 0.67,
  "human_decision_load": 0.54,
  "ai_execution_complexity": 0.38,
  "uncertainty": 0.61
}
```

```bash
estimatorctl estimate --work-item WI-1 --features features.json
```

---

## Resultado

| Salida | Valor | Qué significa |
|---|---|---|
| Human / AI / Verification / Integration Effort | 0.58 / 0.49 / 0.67 / 0.73 | Integration domina porque `integration_complexity` era la más alta (0.83) |
| **Delivery Effort (DEU)** | **6.16** | La predicción global, en la escala propia del framework |
| Confianza | 39% | Reflejo directo de una incertidumbre alta (0.61) |
| Rango (P50 / P80) | 6.16 / 10.51 | "Probablemente 6.16, podría llegar a 10.51" |
| Riesgo | medio | 0.61 cae en la banda media (0.34–0.67) |
| Duración esperada | 1.23 días | Predicción separada del esfuerzo |
| Costo de IA esperado | $0.91 | Predicción económica separada |

---

## Cómo leer estos números

- **DEU no es horas ni Story Points** — no se los presentes así a un stakeholder.
- **Es una predicción, no una promesa** — sirve para compararse después contra la realidad.
- **La confianza es baja a propósito hoy** — todavía no hay historial de calibración.
- **Cada número se puede rastrear** hasta qué dimensión de entrada lo generó — nada es "porque sí".

---

<!-- _class: section -->

# Cerrando el ciclo

---

## De la predicción al aprendizaje

```bash
estimatorctl record-outcome --estimation-id <id> --outcome outcome.json
estimatorctl error-report   --estimation-id <id>
```

- El **Error Report** compara predicho vs. real en cada dimensión (error absoluto, relativo, sesgo).
- Esos datos son exactamente lo que un futuro **Calibration Engine** usará para ajustar los pesos del modelo.
- Un nuevo modelo (`v2`) **nunca sobrescribe** las predicciones hechas con `v1-linear` — siguen siendo reproducibles tal cual fueron.

---

## Lo que ya está construido

| Pieza | Qué hace |
|---|---|
| **Motor de Estimación** | Recibe las 7 features, devuelve la predicción completa, la persiste, acepta el resultado real y calcula el error |
| **Plugin de Claude Code** | Detecta automáticamente cuándo terminó el planning de una feature (Spec-Kit) y dispara todo el flujo sin intervención manual |

**Todavía no**: calibración automática del modelo, y formatos de especificación más allá de Spec Kit (BDD, DDD, TDD, Markdown libre).

---

<!-- _class: lead -->

# Gracias

<p>Repo: github.com/Mirandaberr/delivery-effort-estimator</p>
<p>Manual completo: docs/manual.md</p>
