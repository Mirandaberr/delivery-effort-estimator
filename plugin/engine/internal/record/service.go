package record

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Mirandaberr/delivery-effort-estimator/internal/estimation"
)

// Sentinel errors returned by Service, checkable with errors.Is.
var (
	ErrUnknownEstimation    = errors.New("unknown_estimation_id")
	ErrOutcomeAlreadyExists = errors.New("outcome_already_recorded")
	ErrNoOutcomeRecorded    = errors.New("no_outcome_recorded")
)

// Service orchestrates the deterministic estimation.Predict/Compute logic
// with storage, enforcing the append-only and one-outcome-per-estimation
// rules that the repository interfaces alone don't guarantee at the Go type
// level.
type Service struct {
	estimations EstimationRepository
	outcomes    OutcomeRepository
	now         func() time.Time
}

// NewService wires a Service to its repositories.
func NewService(estimations EstimationRepository, outcomes OutcomeRepository) *Service {
	return &Service{estimations: estimations, outcomes: outcomes, now: time.Now}
}

// EstimateFromJSON parses and clamps a feature vector (FR-001–FR-003) and
// produces + persists a new EstimationRecord.
func (s *Service) EstimateFromJSON(workItemID string, data []byte) (EstimationRecord, error) {
	features, err := estimation.ParseFeatures(data)
	if err != nil {
		return EstimationRecord{}, err
	}
	assumptions := features.Clamp()
	return s.Estimate(workItemID, features, assumptions)
}

// Estimate produces and persists a new EstimationRecord for an already
// validated feature vector. Every call creates a new record — it never
// mutates a prior one, even for the same workItemID (FR-010).
func (s *Service) Estimate(workItemID string, features estimation.EstimationFeatures, assumptions []string) (EstimationRecord, error) {
	id, err := newID()
	if err != nil {
		return EstimationRecord{}, fmt.Errorf("generate id: %w", err)
	}
	if assumptions == nil {
		assumptions = []string{}
	}

	prediction := estimation.Predict(features)
	risks := estimation.DeriveRisks(features)

	rec := EstimationRecord{
		ID:                 id,
		WorkItemID:         workItemID,
		Timestamp:          s.now().UTC(),
		Features:           features,
		Prediction:         prediction,
		Risks:              risks,
		Assumptions:        assumptions,
		ModelVersion:       estimation.ModelVersion,
		CalibrationVersion: CalibrationVersion,
	}

	if err := s.estimations.SaveEstimation(rec); err != nil {
		return EstimationRecord{}, fmt.Errorf("save estimation record: %w", err)
	}
	return rec, nil
}

// GetEstimation fetches a previously persisted, unmodified EstimationRecord.
func (s *Service) GetEstimation(id string) (EstimationRecord, error) {
	rec, ok, err := s.estimations.GetEstimation(id)
	if err != nil {
		return EstimationRecord{}, fmt.Errorf("get estimation record: %w", err)
	}
	if !ok {
		return EstimationRecord{}, ErrUnknownEstimation
	}
	return rec, nil
}

// RecordOutcomeFromJSON parses an OutcomeRecord body and records it.
func (s *Service) RecordOutcomeFromJSON(estimationID string, data []byte) (OutcomeRecord, error) {
	var outcome OutcomeRecord
	if err := json.Unmarshal(data, &outcome); err != nil {
		return OutcomeRecord{}, fmt.Errorf("%w: %v", estimation.ErrInvalidJSON, err)
	}
	return s.RecordOutcome(estimationID, outcome)
}

// RecordOutcome links outcome to an existing EstimationRecord, rejecting an
// unknown estimation id (FR-012) or a second outcome for one that already has
// one (FR-013).
func (s *Service) RecordOutcome(estimationID string, outcome OutcomeRecord) (OutcomeRecord, error) {
	if _, ok, err := s.estimations.GetEstimation(estimationID); err != nil {
		return OutcomeRecord{}, fmt.Errorf("lookup estimation record: %w", err)
	} else if !ok {
		return OutcomeRecord{}, ErrUnknownEstimation
	}

	if _, ok, err := s.outcomes.GetOutcome(estimationID); err != nil {
		return OutcomeRecord{}, fmt.Errorf("lookup outcome record: %w", err)
	} else if ok {
		return OutcomeRecord{}, ErrOutcomeAlreadyExists
	}

	id, err := newID()
	if err != nil {
		return OutcomeRecord{}, fmt.Errorf("generate id: %w", err)
	}
	outcome.ID = id
	outcome.EstimationID = estimationID

	if err := s.outcomes.SaveOutcome(outcome); err != nil {
		return OutcomeRecord{}, fmt.Errorf("save outcome record: %w", err)
	}
	return outcome, nil
}

// ErrorReport computes the comparison between a stored EstimationRecord and
// its OutcomeRecord. It is pure computation — nothing is persisted (FR-015).
func (s *Service) ErrorReport(estimationID string) (estimation.ErrorReport, error) {
	rec, ok, err := s.estimations.GetEstimation(estimationID)
	if err != nil {
		return estimation.ErrorReport{}, fmt.Errorf("lookup estimation record: %w", err)
	}
	if !ok {
		return estimation.ErrorReport{}, ErrUnknownEstimation
	}

	outcome, ok, err := s.outcomes.GetOutcome(estimationID)
	if err != nil {
		return estimation.ErrorReport{}, fmt.Errorf("lookup outcome record: %w", err)
	}
	if !ok {
		return estimation.ErrorReport{}, ErrNoOutcomeRecorded
	}

	actual := estimation.ActualOutcome{
		ActualHumanEffort:        outcome.ActualHumanEffort,
		ActualAIUsage:            outcome.ActualAIUsage,
		ActualAICostUSD:          outcome.ActualAICostUSD,
		ActualLeadTimeDays:       outcome.ActualLeadTimeDays,
		ActualVerificationEffort: outcome.ActualVerificationEffort,
		ActualIntegrationEffort:  outcome.ActualIntegrationEffort,
	}

	return estimation.Compute(estimationID, rec.Prediction, actual), nil
}

// newID generates a random UUID-v4-like identifier without pulling in an
// external dependency beyond the sole storage driver this project already
// needs.
func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
