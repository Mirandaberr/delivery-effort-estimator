package record

import (
	"time"

	"github.com/Mirandaberr/delivery-effort-estimator/internal/estimation"
)

// CalibrationVersion is the calibration state applied to every prediction
// until a Calibration Engine (out of scope for this feature) exists.
const CalibrationVersion = "uncalibrated"

// EstimationRecord is the persisted, immutable result of one estimation run
// (SDD §14). No code path may update or delete a saved record — see
// EstimationRepository.
type EstimationRecord struct {
	ID                   string                        `json:"id"`
	WorkItemID           string                        `json:"work_item_id"`
	Timestamp            time.Time                     `json:"timestamp"`
	SpecificationVersion *string                       `json:"specification_version"`
	PlanningVersion      *string                       `json:"planning_version"`
	RepositoryRevision   *string                       `json:"repository_revision"`
	Features             estimation.EstimationFeatures `json:"features"`
	Prediction           estimation.Prediction         `json:"prediction"`
	Risks                []string                      `json:"risks"`
	Assumptions          []string                      `json:"assumptions"`
	ModelVersion         string                        `json:"model_version"`
	CalibrationVersion   string                        `json:"calibration_version"`
}

// OutcomeRecord captures real measurements collected after a work item is
// delivered, linked to the EstimationRecord it will be compared against
// (SDD §16). At most one OutcomeRecord may exist per EstimationRecord.
type OutcomeRecord struct {
	ID                       string    `json:"id"`
	EstimationID             string    `json:"estimation_id"`
	ActualHumanEffort        float64   `json:"actual_human_effort"`
	ActualAIUsage            float64   `json:"actual_ai_usage"`
	ActualAICostUSD          float64   `json:"actual_ai_cost_usd"`
	ActualLeadTimeDays       float64   `json:"actual_lead_time_days"`
	ActualVerificationEffort float64   `json:"actual_verification_effort"`
	ActualIntegrationEffort  float64   `json:"actual_integration_effort"`
	Rework                   float64   `json:"rework"`
	Incidents                int       `json:"incidents"`
	CompletionTimestamp      time.Time `json:"completion_timestamp"`
}
