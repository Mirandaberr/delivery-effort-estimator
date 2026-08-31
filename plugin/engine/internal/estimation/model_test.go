package estimation

import (
	"math"
	"testing"
)

const epsilon = 1e-4

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestPredictAllZeroFeatures(t *testing.T) {
	p := Predict(EstimationFeatures{})

	cases := map[string]struct{ got, want float64 }{
		"human_effort":        {p.HumanEffort, 0},
		"ai_effort":           {p.AIEffort, 0},
		"verification_effort": {p.VerificationEffort, 0},
		"integration_effort":  {p.IntegrationEffort, 0},
		"delivery_effort":     {p.DeliveryEffortDEU, 0},
		"expected_duration":   {p.ExpectedDurationDays, 0},
		"ai_cost":             {p.ExpectedAICostUSD, 0},
		"confidence":          {p.Confidence, confidenceCeiling},
		"p50":                 {p.PredictionIntervalP50, 0},
	}
	for name, c := range cases {
		if !almostEqual(c.got, c.want) {
			t.Errorf("%s: got %v, want %v", name, c.got, c.want)
		}
	}
	if p.ExpectedAICostEstimated {
		t.Error("expected ExpectedAICostEstimated=false when ai_effort is exactly 0")
	}
	if p.RiskLabel != "low" {
		t.Errorf("expected risk label %q, got %q", "low", p.RiskLabel)
	}
	// Note: P80 == P50 == 0 here is correct, not a bug — the interval scales
	// multiplicatively with DeliveryEffort, and zero effort legitimately has
	// nothing to be uncertain about. The non-zero worked example below is
	// what verifies P80 > P50 in the general case.
}

func TestPredictAllOneFeatures(t *testing.T) {
	f := EstimationFeatures{
		ContextComplexity:      1,
		DomainComplexity:       1,
		IntegrationComplexity:  1,
		VerificationComplexity: 1,
		HumanDecisionLoad:      1,
		AIExecutionComplexity:  1,
		Uncertainty:            1,
	}
	p := Predict(f)

	cases := map[string]struct{ got, want float64 }{
		"human_effort":        {p.HumanEffort, 1.0},
		"ai_effort":           {p.AIEffort, 1.0},
		"verification_effort": {p.VerificationEffort, 1.0},
		"integration_effort":  {p.IntegrationEffort, 1.0},
		"delivery_effort":     {p.DeliveryEffortDEU, 10.0},
		"confidence":          {p.Confidence, confidenceFloor},
		"expected_duration":   {p.ExpectedDurationDays, 2.0},
		"ai_cost":             {p.ExpectedAICostUSD, 3.0},
		"p50":                 {p.PredictionIntervalP50, 10.0},
		"p80":                 {p.PredictionIntervalP80, 19.625},
	}
	for name, c := range cases {
		if !almostEqual(c.got, c.want) {
			t.Errorf("%s: got %v, want %v", name, c.got, c.want)
		}
	}
	if !p.ExpectedAICostEstimated {
		t.Error("expected ExpectedAICostEstimated=true")
	}
	if p.RiskLabel != "high" {
		t.Errorf("expected risk label %q, got %q", "high", p.RiskLabel)
	}
}

func TestPredictWorkedExampleFromSDD(t *testing.T) {
	f := EstimationFeatures{
		ContextComplexity:      0.72,
		DomainComplexity:       0.41,
		IntegrationComplexity:  0.83,
		VerificationComplexity: 0.67,
		HumanDecisionLoad:      0.54,
		AIExecutionComplexity:  0.38,
		Uncertainty:            0.61,
	}
	p := Predict(f)

	cases := map[string]struct{ got, want float64 }{
		"human_effort":        {p.HumanEffort, 0.5815},
		"ai_effort":           {p.AIEffort, 0.4925},
		"verification_effort": {p.VerificationEffort, 0.667},
		"integration_effort":  {p.IntegrationEffort, 0.734},
		"delivery_effort":     {p.DeliveryEffortDEU, 6.15575},
		"confidence":          {p.Confidence, 0.39},
		"expected_duration":   {p.ExpectedDurationDays, 1.23115},
		"ai_cost":             {p.ExpectedAICostUSD, 0.90951},
		"p50":                 {p.PredictionIntervalP50, 6.15575},
		"p80":                 {p.PredictionIntervalP80, 10.51094},
	}
	for name, c := range cases {
		if !almostEqual(c.got, c.want) {
			t.Errorf("%s: got %v, want %v", name, c.got, c.want)
		}
	}
	if p.RiskLabel != "medium" {
		t.Errorf("expected risk label %q, got %q", "medium", p.RiskLabel)
	}
}

func TestRiskLabelThresholds(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{0.0, "low"},
		{0.33, "low"},
		{0.34, "medium"},
		{0.66, "medium"},
		{0.67, "high"},
		{1.0, "high"},
	}
	for _, tt := range tests {
		if got := riskLabelFor(tt.score); got != tt.want {
			t.Errorf("riskLabelFor(%v) = %q, want %q", tt.score, got, tt.want)
		}
	}
}

func TestDeriveRisksFlagsHighDimensions(t *testing.T) {
	f := EstimationFeatures{IntegrationComplexity: 0.9, ContextComplexity: 0.2}
	risks := DeriveRisks(f)
	if len(risks) != 1 {
		t.Fatalf("expected exactly 1 risk note, got %d: %v", len(risks), risks)
	}
}
