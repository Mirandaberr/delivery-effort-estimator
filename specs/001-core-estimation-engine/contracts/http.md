# Contract: `estimatord` HTTP JSON API

Thin wrapper over the same `internal/record.Service` used by the CLI — no business
logic lives in the HTTP layer. All bodies are `application/json`.

## `POST /work-items/{workItemId}/estimations`

Request body: `EstimationFeatures` (same shape as the CLI's `estimate` input).

- `201 Created` → body: `EstimationRecord` (data-model.md).
- `400 Bad Request` → body: `{"error": "missing_feature", "field": "..."}`.

## `GET /estimations/{estimationId}`

- `200 OK` → body: `EstimationRecord`.
- `404 Not Found` → body: `{"error": "unknown_estimation_id"}`.

## `POST /estimations/{estimationId}/outcome`

Request body: `OutcomeRecord` (same shape as the CLI's `record-outcome` input).

- `201 Created` → body: `OutcomeRecord`.
- `404 Not Found` → `{"error": "unknown_estimation_id"}`.
- `409 Conflict` → `{"error": "outcome_already_recorded"}`.

## `GET /estimations/{estimationId}/error-report`

- `200 OK` → body: `ErrorReport`.
- `404 Not Found` → `{"error": "unknown_estimation_id"}` or
  `{"error": "no_outcome_recorded"}`.

## Notes

- No authentication/authorization is in scope for v1 (single-team internal tool,
  consistent with spec.md Assumptions — no end-user UI).
- No pagination/list endpoints in v1 — only fetch-by-id, since spec.md's user stories
  never require enumerating past estimations. A `GET /work-items/{id}/estimations`
  listing endpoint is a natural, low-risk follow-up but is not required by any FR here.
