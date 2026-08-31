Software Design Document

AI-Native Delivery Effort Estimation Framework

1. Objetivo

Diseñar un framework de estimación de esfuerzo de delivery para equipos de desarrollo que utilizan LLMs y agentes de IA.

El framework debe transformar la información disponible durante la especificación y planificación de un requerimiento en una predicción cuantitativa de Delivery Effort, acompañada de incertidumbre, riesgos y métricas explicativas.

La predicción deberá poder compararse posteriormente con los resultados reales del desarrollo, permitiendo medir el error y calibrar progresivamente el modelo.

El objetivo no es reemplazar Story Points por otra escala arbitraria, sino construir un sistema de estimación basado en evidencia, medible y autocorregible.

⸻

2. Problema

Los modelos tradicionales de estimación asumen una relación relativamente estable entre:

Complejidad → Esfuerzo humano → Tiempo → Costo

El desarrollo asistido por IA rompe parcialmente esta relación.

Una tarea puede tener:

* alta complejidad técnica pero bajo esfuerzo humano;
* bajo esfuerzo de implementación pero alto esfuerzo de validación;
* bajo tiempo de desarrollo pero alto consumo de IA;
* alta incertidumbre aunque el código requerido sea simple;
* alto esfuerzo de integración independientemente de la complejidad del código.

Por lo tanto, la estimación debe representar múltiples dimensiones del delivery.

⸻

3. Principios

3.1 Evidence over intuition

Las estimaciones deben basarse en evidencia disponible:

* especificación;
* arquitectura;
* repositorio;
* histórico;
* dependencias;
* resultados anteriores;
* información proporcionada por el ingeniero.

3.2 Effort ≠ Duration

El esfuerzo requerido y el tiempo transcurrido son métricas diferentes.

3.3 Effort ≠ Cost

El costo económico debe calcularse independientemente del esfuerzo.

3.4 Business Value ≠ Delivery Effort

El valor de una iniciativa no debe alterar artificialmente su esfuerzo estimado.

Business Value puede utilizarse posteriormente para priorización y ROI.

3.5 Human and AI effort are different dimensions

La participación de humanos y agentes debe poder medirse independientemente.

3.6 Estimation is a prediction

Una estimación no representa una verdad ni un compromiso.

Representa una predicción basada en la información disponible en un momento determinado.

3.7 Predictions must be measurable

Toda predicción debe poder compararse posteriormente con resultados reales.

3.8 The model must be calibratable

El sistema debe permitir modificar sus parámetros y metodología basándose en evidencia histórica.

3.9 LLMs are not the estimation model

Los LLMs pueden analizar contexto y producir features, pero el motor cuantitativo debe permanecer desacoplado del proveedor de LLM.

⸻

4. Alcance

El framework deberá permitir:

1. analizar un nuevo requerimiento;
2. explorar el repositorio;
3. conversar con el ingeniero;
4. producir o consumir especificaciones;
5. producir un planning técnico;
6. identificar factores de esfuerzo;
7. identificar riesgos e incertidumbres;
8. generar features de estimación;
9. producir una predicción cuantitativa;
10. almacenar la predicción;
11. recolectar métricas reales;
12. comparar predicción versus resultado;
13. analizar errores;
14. calibrar el modelo.

⸻

5. Flujo principal

Requirement
     │
     ▼
Specification
     │
     ▼
Repository Exploration
     │
     ▼
Engineer ↔ LLM Discovery
     │
     ▼
Planning
     │
     ▼
Estimation Features
     │
     ▼
Estimation Engine
     │
     ▼
Prediction
     │
     ▼
Development
     │
     ▼
Actual Measurements
     │
     ▼
Prediction vs Actual
     │
     ▼
Calibration

⸻

6. Specification Layer

El framework no debe imponer un formato único de especificación.

Debe poder consumir:

* SDD;
* BDD;
* DDD;
* TDD;
* Markdown;
* formatos estructurados futuros.

La primera integración será con GitHub Spec Kit.

La especificación debe responder principalmente:

¿Qué problema estamos resolviendo y qué comportamiento debe producir el sistema?

El framework no debe confundir especificación con estimación.

⸻

7. Planning Layer

A partir de la especificación, el sistema debe construir o analizar un planning técnico.

El planning debe identificar, cuando corresponda:

* componentes afectados;
* archivos/módulos potencialmente afectados;
* servicios;
* bases de datos;
* APIs;
* eventos;
* dependencias;
* infraestructura;
* estrategia de testing;
* estrategia de deployment;
* migraciones;
* tareas;
* incertidumbres;
* decisiones pendientes.

El planning constituye una de las principales fuentes de evidencia para el estimador.

⸻

8. Estimation Dimensions

La primera versión del framework utilizará las siguientes dimensiones.

Human Effort

Trabajo humano requerido.

Subdimensiones:

* understanding;
* decision making;
* implementation;
* review;
* validation.

AI Effort

Esfuerzo de ejecución de agentes.

Subdimensiones:

* context acquisition;
* reasoning complexity;
* generation;
* iteration;
* tool execution.

Verification Effort

Esfuerzo necesario para demostrar que el resultado es correcto.

Considerará:

* testing;
* review;
* regression;
* manual validation;
* assurance complexity.

Integration Effort

Esfuerzo necesario para integrar el cambio.

Considerará:

* componentes afectados;
* dependencias;
* APIs;
* eventos;
* bases de datos;
* infraestructura;
* sistemas externos.

Uncertainty

Incertidumbre existente al momento de estimar.

Categorías:

* technical;
* domain;
* business;
* contextual;
* dependency;
* AI-related.

⸻

9. AI Cost

El costo de IA será una dimensión económica independiente.

Podrá incorporar:

* input tokens;
* output tokens;
* cached tokens;
* modelo utilizado;
* tool execution;
* infraestructura asociada.

Los tokens no serán considerados una unidad universal de esfuerzo.

Su propósito inicial será medir costo y comportamiento operacional.

⸻

10. Estimation Features

El sistema transformará la información recopilada en un conjunto estructurado de features.

Ejemplo conceptual:

{
  "context_complexity": 0.72,
  "domain_complexity": 0.41,
  "integration_complexity": 0.83,
  "verification_complexity": 0.67,
  "human_decision_load": 0.54,
  "ai_execution_complexity": 0.38,
  "uncertainty": 0.61
}

Estos valores no deben ser introducidos manualmente necesariamente.

Podrán derivarse de:

* análisis del repositorio;
* especificación;
* planning;
* conversación LLM-humano;
* histórico;
* herramientas de desarrollo.

⸻

11. Estimation Engine

El Estimation Engine será independiente de:

* Claude Code;
* GitHub Spec Kit;
* cualquier proveedor LLM;
* lenguaje de programación;
* sistema de control de versiones.

Su entrada será un conjunto normalizado de features y contexto histórico.

Su salida deberá incluir:

Delivery Effort
Confidence
Prediction Interval
Human Effort
AI Effort
Verification Effort
Integration Effort
Risk
Expected Duration
Expected AI Cost

La primera versión podrá utilizar un modelo determinístico/simple.

No se debe introducir Machine Learning hasta disponer de suficiente información histórica para justificarlo.

⸻

12. Delivery Effort Unit

El framework utilizará inicialmente una unidad relativa propia:

Delivery Effort Unit (DEU).

DEU no representa horas ni Story Points.

Su significado será aprendido empíricamente mediante observaciones históricas.

La relación entre DEU y métricas externas podrá evolucionar.

Ejemplo:

Estimated:
8 DEU
Expected:
6h human effort
1.5 days lead time
$2 AI cost

La equivalencia anterior es una predicción y no una definición de DEU.

⸻

13. Confidence and Prediction Interval

El sistema no deberá producir únicamente:

Estimated effort = X

Debe producir una estimación probabilística cuando exista suficiente información.

Ejemplo conceptual:

P50 = 6h
P80 = 10h

Esto permite expresar incertidumbre explícitamente.

⸻

14. Estimation Record

Cada estimación deberá quedar registrada.

Conceptualmente:

EstimationRecord
id
workItemId
timestamp
specificationVersion
planningVersion
repositoryRevision
features
prediction
confidence
predictionInterval
risks
assumptions
modelVersion
calibrationVersion

La versión del modelo es obligatoria para poder reconstruir cómo se produjo una predicción histórica.

⸻

15. Actual Measurements

Durante y después del desarrollo se recopilarán métricas reales.

Posibles fuentes:

* Git;
* Pull Requests;
* CI/CD;
* testing;
* deployment;
* agentes;
* LLM APIs;
* IDE;
* observabilidad.

Posibles métricas:

human effort
elapsed time
lead time
AI tokens
AI cost
iterations
rework
files changed
services changed
test execution
deployment attempts
incidents

La primera versión no necesita integrar todas estas fuentes.

El framework debe diseñarse para permitir incorporarlas progresivamente.

⸻

16. Outcome Record

Al finalizar el trabajo:

OutcomeRecord
estimationId
actualHumanEffort
actualAIUsage
actualAICost
actualLeadTime
actualVerificationEffort
actualIntegrationEffort
rework
incidents
completionTimestamp

⸻

17. Calibration

El sistema deberá comparar:

Prediction
     vs
Actual

y calcular:

* absolute error;
* relative error;
* bias;
* mean error;
* prediction interval coverage;
* error por categoría;
* error por proyecto;
* error por tipo de cambio.

La calibración no deberá reaccionar automáticamente ante una única desviación.

Debe utilizar evidencia estadística suficiente para identificar sesgos persistentes.

⸻

18. Retrospective Analysis

El framework debe permitir preguntas como:

¿Dónde estamos subestimando?
¿Dónde estamos sobreestimando?
¿Qué variables explican los errores?
¿Qué tipos de cambios presentan mayor desviación?
¿Qué factores son realmente predictivos?
¿El modelo está mejorando?
¿La estimación es mejor que Story Points?
¿La estimación es mejor que horas?

Esto convierte el sistema en una herramienta experimental además de operacional.

⸻

19. Model Evolution

El modelo debe estar versionado.

Model v1
   ↓
Observations
   ↓
Calibration
   ↓
Model v2

Nunca se deben sobrescribir las estimaciones históricas.

Una estimación realizada con Model v1 debe seguir siendo reproducible bajo Model v1.

⸻

20. Arquitectura conceptual

                  ┌──────────────────────┐
                  │ Specification Tools  │
                  │ Spec Kit / Other     │
                  └──────────┬───────────┘
                             │
                             ▼
                  ┌──────────────────────┐
                  │ Discovery Layer      │
                  │ LLM + Engineer       │
                  └──────────┬───────────┘
                             │
                 ┌───────────┴───────────┐
                 ▼                       ▼
          Repository Analysis        Planning
                 │                       │
                 └───────────┬───────────┘
                             ▼
                  ┌──────────────────────┐
                  │ Feature Extraction   │
                  └──────────┬───────────┘
                             ▼
                  ┌──────────────────────┐
                  │ Estimation Engine    │
                  └──────────┬───────────┘
                             ▼
                       Prediction
                             │
                             ▼
                       Development
                             │
                    ┌────────┴────────┐
                    ▼                 ▼
                 Telemetry          Git/CI
                    │                 │
                    └────────┬────────┘
                             ▼
                  ┌──────────────────────┐
                  │ Outcome Measurement  │
                  └──────────┬───────────┘
                             ▼
                  ┌──────────────────────┐
                  │ Calibration Engine   │
                  └──────────┬───────────┘
                             │
                             ▼
                       Model Version

⸻

21. Integración inicial

La primera implementación deberá integrarse con:

* Claude Code;
* GitHub Spec Kit;
* Git;
* repositorios locales.

La integración con otros agentes debe ser posible sin modificar el núcleo del framework.

⸻

22. No Goals

La primera versión no pretende:

* reemplazar Jira;
* reemplazar Scrum;
* reemplazar Product Management;
* eliminar Story Points;
* determinar Business Value;
* controlar la ejecución de agentes;
* convertirse en un sistema completo de observabilidad;
* utilizar Machine Learning desde el primer día.

⸻

23. Criterio de éxito

El framework será considerado exitoso si demuestra progresivamente que:

1. puede producir estimaciones reproducibles;
2. puede explicar los factores que determinan una estimación;
3. puede medir Prediction vs Actual;
4. puede identificar sesgos;
5. puede mejorar sus predicciones mediante calibración;
6. puede correlacionar Delivery Effort con métricas reales;
7. puede demostrar empíricamente si su capacidad predictiva supera métodos tradicionales.

La superioridad respecto de Story Points, horas u otros mecanismos no será asumida; deberá demostrarse mediante datos.
