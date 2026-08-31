package estimation

import "testing"

func TestComputeDimensionErrors(t *testing.T) {
	prediction := Prediction{
		HumanEffort:          5.0,
		AIEffort:             0,
		VerificationEffort:   3.0,
		IntegrationEffort:    2.0,
		ExpectedAICostUSD:    1.0,
		ExpectedDurationDays: 2.0,
	}
	actual := ActualOutcome{
		ActualHumanEffort:        6.0,
		ActualAIUsage:            100,
		ActualVerificationEffort: 3.0,
		ActualIntegrationEffort:  1.0,
		ActualAICostUSD:          1.5,
		ActualLeadTimeDays:       2.5,
	}

	report := Compute("est-1", prediction, actual)
	if report.EstimationID != "est-1" {
		t.Errorf("expected estimation id %q, got %q", "est-1", report.EstimationID)
	}
	if len(report.Dimensions) != 6 {
		t.Fatalf("expected 6 dimensions, got %d", len(report.Dimensions))
	}

	byName := map[string]DimensionError{}
	for _, d := range report.Dimensions {
		byName[d.Name] = d
	}

	human := byName["human_effort"]
	if !almostEqual(human.AbsoluteError, 1.0) || !almostEqual(human.Bias, 1.0) {
		t.Errorf("human_effort: got abs=%v bias=%v", human.AbsoluteError, human.Bias)
	}
	if human.RelativeError == nil || !almostEqual(*human.RelativeError, 0.2) {
		t.Errorf("human_effort relative error: got %v, want 0.2", human.RelativeError)
	}

	aiEffort := byName["ai_effort"]
	if aiEffort.RelativeError != nil {
		t.Errorf("ai_effort: expected nil relative error when predicted is 0, got %v", *aiEffort.RelativeError)
	}
	if !almostEqual(aiEffort.Bias, 100) {
		t.Errorf("ai_effort bias: got %v, want 100", aiEffort.Bias)
	}

	integration := byName["integration_effort"]
	if !almostEqual(integration.Bias, -1.0) {
		t.Errorf("integration_effort bias: got %v, want -1.0 (overestimated)", integration.Bias)
	}

	verification := byName["verification_effort"]
	if !almostEqual(verification.AbsoluteError, 0) {
		t.Errorf("verification_effort: expected zero error for an exact match, got %v", verification.AbsoluteError)
	}
}
