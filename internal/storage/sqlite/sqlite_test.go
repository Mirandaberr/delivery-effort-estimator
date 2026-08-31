package sqlite

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Mirandaberr/delivery-effort-estimator/internal/estimation"
	"github.com/Mirandaberr/delivery-effort-estimator/internal/record"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "estimator.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func minimalEstimationRecord(id string) record.EstimationRecord {
	features := estimation.EstimationFeatures{
		ContextComplexity: 0.5, DomainComplexity: 0.4, IntegrationComplexity: 0.3,
		VerificationComplexity: 0.6, HumanDecisionLoad: 0.2, AIExecutionComplexity: 0.1,
		Uncertainty: 0.3,
	}
	return record.EstimationRecord{
		ID:                 id,
		WorkItemID:         "WI-1",
		Timestamp:          time.Now().UTC().Truncate(time.Second),
		Features:           features,
		Prediction:         estimation.Predict(features),
		Risks:              []string{},
		Assumptions:        []string{},
		ModelVersion:       estimation.ModelVersion,
		CalibrationVersion: "uncalibrated",
	}
}

func TestSaveAndGetEstimationRoundTrip(t *testing.T) {
	store := openTestStore(t)
	rec := minimalEstimationRecord("est-1")

	if err := store.SaveEstimation(rec); err != nil {
		t.Fatalf("SaveEstimation: %v", err)
	}

	got, ok, err := store.GetEstimation("est-1")
	if err != nil {
		t.Fatalf("GetEstimation: %v", err)
	}
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Prediction != rec.Prediction {
		t.Errorf("prediction mismatch: got %+v want %+v", got.Prediction, rec.Prediction)
	}
	if got.Features != rec.Features {
		t.Errorf("features mismatch: got %+v want %+v", got.Features, rec.Features)
	}
	if got.WorkItemID != rec.WorkItemID {
		t.Errorf("work item id mismatch: got %s want %s", got.WorkItemID, rec.WorkItemID)
	}
	if !got.Timestamp.Equal(rec.Timestamp) {
		t.Errorf("timestamp mismatch: got %v want %v", got.Timestamp, rec.Timestamp)
	}
}

func TestGetEstimationUnknownIDReturnsNotFound(t *testing.T) {
	store := openTestStore(t)
	_, ok, err := store.GetEstimation("does-not-exist")
	if err != nil {
		t.Fatalf("GetEstimation: %v", err)
	}
	if ok {
		t.Fatal("expected not found for unknown id")
	}
}

func TestDuplicateOutcomeRejectedAtDBLevel(t *testing.T) {
	store := openTestStore(t)
	rec := minimalEstimationRecord("est-2")
	if err := store.SaveEstimation(rec); err != nil {
		t.Fatalf("SaveEstimation: %v", err)
	}

	outcome := record.OutcomeRecord{
		ID:                  "out-1",
		EstimationID:        "est-2",
		CompletionTimestamp: time.Now().UTC().Truncate(time.Second),
	}
	if err := store.SaveOutcome(outcome); err != nil {
		t.Fatalf("first SaveOutcome: %v", err)
	}

	outcome.ID = "out-2"
	if err := store.SaveOutcome(outcome); err == nil {
		t.Fatal("expected second SaveOutcome for the same estimation_id to fail (UNIQUE constraint)")
	}
}

func TestGetOutcomeUnknownEstimationReturnsNotFound(t *testing.T) {
	store := openTestStore(t)
	_, ok, err := store.GetOutcome("does-not-exist")
	if err != nil {
		t.Fatalf("GetOutcome: %v", err)
	}
	if ok {
		t.Fatal("expected not found for unknown estimation id")
	}
}
