package estimation

import (
	"errors"
	"testing"
)

func TestParseFeaturesMissingField(t *testing.T) {
	data := []byte(`{
		"context_complexity": 0.5,
		"domain_complexity": 0.5,
		"integration_complexity": 0.5,
		"verification_complexity": 0.5,
		"human_decision_load": 0.5,
		"ai_execution_complexity": 0.5
	}`)

	_, err := ParseFeatures(data)
	field, ok := IsMissingFeatureError(err)
	if !ok {
		t.Fatalf("expected MissingFeatureError, got %v", err)
	}
	if field != "uncertainty" {
		t.Errorf("expected missing field %q, got %q", "uncertainty", field)
	}
}

func TestParseFeaturesInvalidJSON(t *testing.T) {
	_, err := ParseFeatures([]byte(`not json`))
	if !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestParseFeaturesCompleteVector(t *testing.T) {
	data := []byte(`{
		"context_complexity": 0.72,
		"domain_complexity": 0.41,
		"integration_complexity": 0.83,
		"verification_complexity": 0.67,
		"human_decision_load": 0.54,
		"ai_execution_complexity": 0.38,
		"uncertainty": 0.61
	}`)

	got, err := ParseFeatures(data)
	if err != nil {
		t.Fatalf("ParseFeatures: %v", err)
	}
	want := EstimationFeatures{
		ContextComplexity:      0.72,
		DomainComplexity:       0.41,
		IntegrationComplexity:  0.83,
		VerificationComplexity: 0.67,
		HumanDecisionLoad:      0.54,
		AIExecutionComplexity:  0.38,
		Uncertainty:            0.61,
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestClampOutOfRangeValues(t *testing.T) {
	f := EstimationFeatures{
		ContextComplexity: 1.4,
		Uncertainty:       -0.2,
	}
	assumptions := f.Clamp()

	if f.ContextComplexity != 1.0 {
		t.Errorf("expected context_complexity clamped to 1.0, got %v", f.ContextComplexity)
	}
	if f.Uncertainty != 0.0 {
		t.Errorf("expected uncertainty clamped to 0.0, got %v", f.Uncertainty)
	}
	if len(assumptions) != 2 {
		t.Errorf("expected 2 assumption notes, got %d: %v", len(assumptions), assumptions)
	}
}

func TestClampWithinRangeProducesNoAssumptions(t *testing.T) {
	f := EstimationFeatures{ContextComplexity: 0.5, Uncertainty: 0.5}
	assumptions := f.Clamp()
	if len(assumptions) != 0 {
		t.Errorf("expected no assumptions, got %v", assumptions)
	}
}
