package estimation

import (
	"fmt"
	"math"
)

// ModelVersion identifies the deterministic model logic used to produce a
// Prediction. It is stored on every EstimationRecord so a historical
// prediction stays reproducible after later model changes (Constitution
// Principle VIII).
const ModelVersion = "v1-linear"

const (
	confidenceFloor     = 0.05
	confidenceCeiling   = 0.70
	intervalFloorSpread = 0.25
	durationFactor      = 0.2
	aiCostFactor        = 0.3

	riskLowUpperBound    = 0.34
	riskMediumUpperBound = 0.67

	riskFeatureThreshold = 0.7
)

// Prediction is the structured output of one deterministic estimation run
// (SDD §11).
type Prediction struct {
	DeliveryEffortDEU       float64 `json:"delivery_effort_deu"`
	Confidence              float64 `json:"confidence"`
	PredictionIntervalP50   float64 `json:"prediction_interval_p50"`
	PredictionIntervalP80   float64 `json:"prediction_interval_p80"`
	HumanEffort             float64 `json:"human_effort"`
	AIEffort                float64 `json:"ai_effort"`
	VerificationEffort      float64 `json:"verification_effort"`
	IntegrationEffort       float64 `json:"integration_effort"`
	RiskScore               float64 `json:"risk_score"`
	RiskLabel               string  `json:"risk_label"`
	ExpectedDurationDays    float64 `json:"expected_duration_days"`
	ExpectedAICostUSD       float64 `json:"expected_ai_cost_usd"`
	ExpectedAICostEstimated bool    `json:"expected_ai_cost_estimated"`
}

// Predict computes a deterministic Prediction from a feature vector. Same
// input + same ModelVersion always yields the same output (FR-005) — no I/O,
// no LLM call, no randomness (Constitution Principle IX).
func Predict(f EstimationFeatures) Prediction {
	humanEffort := 0.45*f.HumanDecisionLoad + 0.30*f.ContextComplexity + 0.15*f.DomainComplexity + 0.10*f.Uncertainty
	aiEffort := 0.50*f.AIExecutionComplexity + 0.25*f.ContextComplexity + 0.15*f.DomainComplexity + 0.10*f.Uncertainty
	verificationEffort := 0.55*f.VerificationComplexity + 0.20*f.IntegrationComplexity + 0.15*f.Uncertainty + 0.10*f.DomainComplexity
	integrationEffort := 0.70*f.IntegrationComplexity + 0.15*f.DomainComplexity + 0.15*f.Uncertainty

	deliveryEffort := 10 * (0.35*humanEffort + 0.20*aiEffort + 0.25*verificationEffort + 0.20*integrationEffort)

	// Confidence is capped below 1.0 even at zero uncertainty: with no
	// calibration history yet (calibration_version == "uncalibrated"), the
	// model must not claim high confidence regardless of input uncertainty.
	confidence := clampRange(1-f.Uncertainty, confidenceFloor, confidenceCeiling)

	// The interval always expresses some uncertainty (spread never reaches
	// 0), satisfying the edge case that a Prediction Interval must exist even
	// in the best-confidence scenario.
	spread := intervalFloorSpread + 0.75*(1-confidence)
	p50 := deliveryEffort
	p80 := deliveryEffort * (1 + spread)

	riskScore := f.Uncertainty
	riskLabel := riskLabelFor(riskScore)

	expectedDuration := durationFactor * deliveryEffort

	// A zero AIEffort is distinguished from "cost could not be estimated" via
	// ExpectedAICostEstimated, rather than an ambiguous bare 0.
	costEstimated := aiEffort > 0
	var aiCost float64
	if costEstimated {
		aiCost = aiCostFactor * aiEffort * deliveryEffort
	}

	return Prediction{
		DeliveryEffortDEU:       deliveryEffort,
		Confidence:              confidence,
		PredictionIntervalP50:   p50,
		PredictionIntervalP80:   p80,
		HumanEffort:             humanEffort,
		AIEffort:                aiEffort,
		VerificationEffort:      verificationEffort,
		IntegrationEffort:       integrationEffort,
		RiskScore:               riskScore,
		RiskLabel:               riskLabel,
		ExpectedDurationDays:    expectedDuration,
		ExpectedAICostUSD:       aiCost,
		ExpectedAICostEstimated: costEstimated,
	}
}

// DeriveRisks returns human-readable notes for any input dimension that is
// itself high enough to be worth flagging directly on the EstimationRecord.
func DeriveRisks(f EstimationFeatures) []string {
	risks := []string{}
	check := func(name string, v float64) {
		if v > riskFeatureThreshold {
			risks = append(risks, fmt.Sprintf("%s is high (%.2f)", name, v))
		}
	}
	check("context_complexity", f.ContextComplexity)
	check("domain_complexity", f.DomainComplexity)
	check("integration_complexity", f.IntegrationComplexity)
	check("verification_complexity", f.VerificationComplexity)
	check("human_decision_load", f.HumanDecisionLoad)
	check("ai_execution_complexity", f.AIExecutionComplexity)
	check("uncertainty", f.Uncertainty)
	return risks
}

func riskLabelFor(score float64) string {
	switch {
	case score < riskLowUpperBound:
		return "low"
	case score < riskMediumUpperBound:
		return "medium"
	default:
		return "high"
	}
}

func clampRange(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}
