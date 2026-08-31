package record

import (
	"errors"
	"testing"

	"github.com/jmirandev/delivery-effort-estimator/internal/estimation"
)

func TestEstimateCreatesDistinctRecordsPerCall(t *testing.T) {
	svc := NewService(newFakeEstimationRepo(), newFakeOutcomeRepo())
	features := estimation.EstimationFeatures{Uncertainty: 0.5}

	first, err := svc.Estimate("WI-1", features, nil)
	if err != nil {
		t.Fatalf("first Estimate: %v", err)
	}
	second, err := svc.Estimate("WI-1", features, nil)
	if err != nil {
		t.Fatalf("second Estimate: %v", err)
	}

	if first.ID == second.ID {
		t.Fatal("expected distinct ids for repeated estimation of the same work item")
	}

	gotFirst, err := svc.GetEstimation(first.ID)
	if err != nil {
		t.Fatalf("GetEstimation(first): %v", err)
	}
	if gotFirst.Prediction != first.Prediction {
		t.Error("first record's prediction changed after a second estimation was made")
	}
}

func TestGetEstimationUnknownID(t *testing.T) {
	svc := NewService(newFakeEstimationRepo(), newFakeOutcomeRepo())
	_, err := svc.GetEstimation("does-not-exist")
	if !errors.Is(err, ErrUnknownEstimation) {
		t.Fatalf("expected ErrUnknownEstimation, got %v", err)
	}
}

func TestRecordOutcomeRejectsUnknownEstimation(t *testing.T) {
	svc := NewService(newFakeEstimationRepo(), newFakeOutcomeRepo())
	_, err := svc.RecordOutcome("does-not-exist", OutcomeRecord{})
	if !errors.Is(err, ErrUnknownEstimation) {
		t.Fatalf("expected ErrUnknownEstimation, got %v", err)
	}
}

func TestRecordOutcomeRejectsDuplicate(t *testing.T) {
	svc := NewService(newFakeEstimationRepo(), newFakeOutcomeRepo())
	rec, err := svc.Estimate("WI-1", estimation.EstimationFeatures{}, nil)
	if err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	if _, err := svc.RecordOutcome(rec.ID, OutcomeRecord{}); err != nil {
		t.Fatalf("first RecordOutcome: %v", err)
	}
	if _, err := svc.RecordOutcome(rec.ID, OutcomeRecord{}); !errors.Is(err, ErrOutcomeAlreadyExists) {
		t.Fatalf("expected ErrOutcomeAlreadyExists, got %v", err)
	}
}

func TestErrorReportRequiresOutcomeFirst(t *testing.T) {
	svc := NewService(newFakeEstimationRepo(), newFakeOutcomeRepo())
	rec, err := svc.Estimate("WI-1", estimation.EstimationFeatures{}, nil)
	if err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	if _, err := svc.ErrorReport(rec.ID); !errors.Is(err, ErrNoOutcomeRecorded) {
		t.Fatalf("expected ErrNoOutcomeRecorded, got %v", err)
	}

	if _, err := svc.RecordOutcome(rec.ID, OutcomeRecord{ActualHumanEffort: 1}); err != nil {
		t.Fatalf("RecordOutcome: %v", err)
	}

	report, err := svc.ErrorReport(rec.ID)
	if err != nil {
		t.Fatalf("ErrorReport: %v", err)
	}
	if report.EstimationID != rec.ID {
		t.Errorf("expected estimation id %q, got %q", rec.ID, report.EstimationID)
	}
}

func TestErrorReportUnknownEstimation(t *testing.T) {
	svc := NewService(newFakeEstimationRepo(), newFakeOutcomeRepo())
	if _, err := svc.ErrorReport("does-not-exist"); !errors.Is(err, ErrUnknownEstimation) {
		t.Fatalf("expected ErrUnknownEstimation, got %v", err)
	}
}
