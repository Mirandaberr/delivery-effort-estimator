// Command estimatord is the HTTP JSON contract described in
// specs/001-core-estimation-engine/contracts/http.md — a thin wrapper over
// internal/record.Service; no business logic lives here.
package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jmirandev/delivery-effort-estimator/internal/estimation"
	"github.com/jmirandev/delivery-effort-estimator/internal/record"
	"github.com/jmirandev/delivery-effort-estimator/internal/storage/sqlite"
)

const defaultDBPath = "./data/estimator.db"

func main() {
	dbPath := os.Getenv("ESTIMATOR_DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	store, err := sqlite.Open(dbPath)
	if err != nil {
		log.Fatalf("open storage: %v", err)
	}
	defer store.Close()

	svc := record.NewService(store, store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /work-items/{workItemId}/estimations", handleEstimate(svc))
	mux.HandleFunc("GET /estimations/{estimationId}", handleGetEstimation(svc))
	mux.HandleFunc("POST /estimations/{estimationId}/outcome", handleRecordOutcome(svc))
	mux.HandleFunc("GET /estimations/{estimationId}/error-report", handleErrorReport(svc))

	addr := os.Getenv("ESTIMATOR_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("estimatord listening on %s (db: %s)", addr, dbPath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleEstimate(svc *record.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workItemID := r.PathValue("workItemId")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", map[string]string{"detail": err.Error()})
			return
		}

		rec, err := svc.EstimateFromJSON(workItemID, data)
		if err != nil {
			if field, ok := estimation.IsMissingFeatureError(err); ok {
				writeError(w, http.StatusBadRequest, "missing_feature", map[string]string{"field": field})
				return
			}
			if errors.Is(err, estimation.ErrInvalidJSON) {
				writeError(w, http.StatusBadRequest, "invalid_json", map[string]string{"detail": err.Error()})
				return
			}
			writeError(w, http.StatusInternalServerError, "internal_error", map[string]string{"detail": err.Error()})
			return
		}

		writeJSON(w, http.StatusCreated, rec)
	}
}

func handleGetEstimation(svc *record.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("estimationId")
		rec, err := svc.GetEstimation(id)
		if err != nil {
			if errors.Is(err, record.ErrUnknownEstimation) {
				writeError(w, http.StatusNotFound, "unknown_estimation_id", map[string]string{"estimation_id": id})
				return
			}
			writeError(w, http.StatusInternalServerError, "internal_error", map[string]string{"detail": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, rec)
	}
}

func handleRecordOutcome(svc *record.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("estimationId")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", map[string]string{"detail": err.Error()})
			return
		}

		outcome, err := svc.RecordOutcomeFromJSON(id, data)
		if err != nil {
			switch {
			case errors.Is(err, record.ErrUnknownEstimation):
				writeError(w, http.StatusNotFound, "unknown_estimation_id", map[string]string{"estimation_id": id})
			case errors.Is(err, record.ErrOutcomeAlreadyExists):
				writeError(w, http.StatusConflict, "outcome_already_recorded", map[string]string{"estimation_id": id})
			case errors.Is(err, estimation.ErrInvalidJSON):
				writeError(w, http.StatusBadRequest, "invalid_json", map[string]string{"detail": err.Error()})
			default:
				writeError(w, http.StatusInternalServerError, "internal_error", map[string]string{"detail": err.Error()})
			}
			return
		}

		writeJSON(w, http.StatusCreated, outcome)
	}
}

func handleErrorReport(svc *record.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("estimationId")
		report, err := svc.ErrorReport(id)
		if err != nil {
			switch {
			case errors.Is(err, record.ErrUnknownEstimation):
				writeError(w, http.StatusNotFound, "unknown_estimation_id", map[string]string{"estimation_id": id})
			case errors.Is(err, record.ErrNoOutcomeRecorded):
				writeError(w, http.StatusNotFound, "no_outcome_recorded", map[string]string{"estimation_id": id})
			default:
				writeError(w, http.StatusInternalServerError, "internal_error", map[string]string{"detail": err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusOK, report)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code string, extra map[string]string) {
	payload := map[string]interface{}{"error": code}
	for k, v := range extra {
		payload[k] = v
	}
	writeJSON(w, status, payload)
}
