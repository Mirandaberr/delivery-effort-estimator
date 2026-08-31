package record

// EstimationRepository persists EstimationRecords. It deliberately exposes no
// update or delete method: an estimation, once saved, is immutable
// (Constitution Principle VIII) — re-estimating a work item calls
// SaveEstimation again with a new id (FR-010).
type EstimationRepository interface {
	SaveEstimation(EstimationRecord) error
	GetEstimation(id string) (EstimationRecord, bool, error)
}

// OutcomeRepository persists OutcomeRecords, at most one per estimation id.
// Like EstimationRepository, it exposes no update/delete method.
type OutcomeRepository interface {
	SaveOutcome(OutcomeRecord) error
	GetOutcome(estimationID string) (OutcomeRecord, bool, error)
}
