package estimation

import "math"

// ActualOutcome carries the real-world measurements needed to compare against
// a stored Prediction (SDD §16, OutcomeRecord).
type ActualOutcome struct {
	ActualHumanEffort        float64
	ActualAIUsage            float64
	ActualAICostUSD          float64
	ActualLeadTimeDays       float64
	ActualVerificationEffort float64
	ActualIntegrationEffort  float64
}

// DimensionError is the computed comparison for one dimension present in
// both a Prediction and an ActualOutcome.
type DimensionError struct {
	Name          string   `json:"name"`
	Predicted     float64  `json:"predicted"`
	Actual        float64  `json:"actual"`
	AbsoluteError float64  `json:"absolute_error"`
	RelativeError *float64 `json:"relative_error"`
	Bias          float64  `json:"bias"`
}

// ErrorReport is the on-demand, non-persisted comparison between one
// EstimationRecord's Prediction and its OutcomeRecord (FR-014, FR-015).
type ErrorReport struct {
	EstimationID string           `json:"estimation_id"`
	Dimensions   []DimensionError `json:"dimensions"`
}

// Compute returns the absolute error, relative error, and bias for every
// dimension present in both p and actual. Raw numeric deltas are used
// deliberately — establishing the DEU-to-real-unit mapping is the future
// Calibration Engine's job (SDD §12, §17), not this engine's (see
// research.md "Cross-dimension error comparison").
func Compute(estimationID string, p Prediction, actual ActualOutcome) ErrorReport {
	dims := []struct {
		name      string
		predicted float64
		actual    float64
	}{
		{"human_effort", p.HumanEffort, actual.ActualHumanEffort},
		{"ai_effort", p.AIEffort, actual.ActualAIUsage},
		{"verification_effort", p.VerificationEffort, actual.ActualVerificationEffort},
		{"integration_effort", p.IntegrationEffort, actual.ActualIntegrationEffort},
		{"ai_cost_usd", p.ExpectedAICostUSD, actual.ActualAICostUSD},
		{"lead_time_days", p.ExpectedDurationDays, actual.ActualLeadTimeDays},
	}

	report := ErrorReport{
		EstimationID: estimationID,
		Dimensions:   make([]DimensionError, 0, len(dims)),
	}

	for _, d := range dims {
		absErr := math.Abs(d.actual - d.predicted)
		bias := d.actual - d.predicted

		var relErr *float64
		if d.predicted != 0 {
			r := absErr / d.predicted
			relErr = &r
		}

		report.Dimensions = append(report.Dimensions, DimensionError{
			Name:          d.name,
			Predicted:     d.predicted,
			Actual:        d.actual,
			AbsoluteError: absErr,
			RelativeError: relErr,
			Bias:          bias,
		})
	}

	return report
}
