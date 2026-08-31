package estimation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EstimationFeatures is the normalized [0,1] input vector for one estimation
// run (SDD §10).
type EstimationFeatures struct {
	ContextComplexity      float64 `json:"context_complexity"`
	DomainComplexity       float64 `json:"domain_complexity"`
	IntegrationComplexity  float64 `json:"integration_complexity"`
	VerificationComplexity float64 `json:"verification_complexity"`
	HumanDecisionLoad      float64 `json:"human_decision_load"`
	AIExecutionComplexity  float64 `json:"ai_execution_complexity"`
	Uncertainty            float64 `json:"uncertainty"`
}

var featureFieldOrder = []string{
	"context_complexity",
	"domain_complexity",
	"integration_complexity",
	"verification_complexity",
	"human_decision_load",
	"ai_execution_complexity",
	"uncertainty",
}

// ErrInvalidJSON marks a request body that could not be parsed as JSON at all,
// as distinct from a well-formed body missing a required dimension.
var ErrInvalidJSON = errors.New("invalid_json")

// MissingFeatureError reports which required dimension (FR-001) was absent
// from a request (FR-002).
type MissingFeatureError struct {
	Field string
}

func (e *MissingFeatureError) Error() string {
	return fmt.Sprintf("missing required feature: %s", e.Field)
}

// IsMissingFeatureError extracts the missing field name from err, if err (or
// something it wraps) is a *MissingFeatureError.
func IsMissingFeatureError(err error) (string, bool) {
	var mfe *MissingFeatureError
	if errors.As(err, &mfe) {
		return mfe.Field, true
	}
	return "", false
}

// ParseFeatures decodes a JSON feature vector, requiring all seven dimensions
// to be present (FR-002). Out-of-range clamping is a separate step: call
// Clamp on the result.
func ParseFeatures(data []byte) (EstimationFeatures, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return EstimationFeatures{}, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	values := make(map[string]float64, len(featureFieldOrder))
	for _, field := range featureFieldOrder {
		rawValue, ok := raw[field]
		if !ok {
			return EstimationFeatures{}, &MissingFeatureError{Field: field}
		}
		var v float64
		if err := json.Unmarshal(rawValue, &v); err != nil {
			return EstimationFeatures{}, fmt.Errorf("%w: field %s: %v", ErrInvalidJSON, field, err)
		}
		values[field] = v
	}

	return EstimationFeatures{
		ContextComplexity:      values["context_complexity"],
		DomainComplexity:       values["domain_complexity"],
		IntegrationComplexity:  values["integration_complexity"],
		VerificationComplexity: values["verification_complexity"],
		HumanDecisionLoad:      values["human_decision_load"],
		AIExecutionComplexity:  values["ai_execution_complexity"],
		Uncertainty:            values["uncertainty"],
	}, nil
}

// Clamp bounds every dimension to [0,1] in place, returning a human-readable
// adjustment note for each value that was out of range (FR-003).
func (f *EstimationFeatures) Clamp() []string {
	assumptions := []string{}
	clamp := func(name string, v *float64) {
		if *v < 0 || *v > 1 {
			original := *v
			if *v < 0 {
				*v = 0
			} else {
				*v = 1
			}
			assumptions = append(assumptions, fmt.Sprintf("clamped %s from %.2f to %.2f", name, original, *v))
		}
	}
	clamp("context_complexity", &f.ContextComplexity)
	clamp("domain_complexity", &f.DomainComplexity)
	clamp("integration_complexity", &f.IntegrationComplexity)
	clamp("verification_complexity", &f.VerificationComplexity)
	clamp("human_decision_load", &f.HumanDecisionLoad)
	clamp("ai_execution_complexity", &f.AIExecutionComplexity)
	clamp("uncertainty", &f.Uncertainty)
	return assumptions
}
