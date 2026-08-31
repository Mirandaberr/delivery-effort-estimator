// Package sqlite implements internal/record's repository interfaces on top
// of an embedded SQLite file, using the pure-Go modernc.org/sqlite driver so
// the resulting binary needs no C toolchain (research.md).
package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/Mirandaberr/delivery-effort-estimator/internal/record"
)

// Store implements both record.EstimationRepository and
// record.OutcomeRepository against a single SQLite file.
type Store struct {
	db *sql.DB
}

var (
	_ record.EstimationRepository = (*Store)(nil)
	_ record.OutcomeRepository    = (*Store)(nil)
)

// Open creates the parent directory for path if needed, opens (or creates)
// the SQLite file there, and applies the schema migration.
func Open(path string) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	return &Store{db: db}, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error {
	return s.db.Close()
}

// SaveEstimation inserts a new estimation_records row. There is no update
// path — see internal/record.EstimationRepository.
func (s *Store) SaveEstimation(rec record.EstimationRecord) error {
	featuresJSON, err := json.Marshal(rec.Features)
	if err != nil {
		return fmt.Errorf("marshal features: %w", err)
	}
	predictionJSON, err := json.Marshal(rec.Prediction)
	if err != nil {
		return fmt.Errorf("marshal prediction: %w", err)
	}
	risksJSON, err := json.Marshal(rec.Risks)
	if err != nil {
		return fmt.Errorf("marshal risks: %w", err)
	}
	assumptionsJSON, err := json.Marshal(rec.Assumptions)
	if err != nil {
		return fmt.Errorf("marshal assumptions: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO estimation_records (
			id, work_item_id, timestamp, specification_version, planning_version,
			repository_revision, features_json, prediction_json, confidence,
			prediction_interval_p50, prediction_interval_p80, risks_json,
			assumptions_json, model_version, calibration_version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.WorkItemID, rec.Timestamp.Format(time.RFC3339Nano),
		rec.SpecificationVersion, rec.PlanningVersion, rec.RepositoryRevision,
		string(featuresJSON), string(predictionJSON), rec.Prediction.Confidence,
		rec.Prediction.PredictionIntervalP50, rec.Prediction.PredictionIntervalP80,
		string(risksJSON), string(assumptionsJSON), rec.ModelVersion, rec.CalibrationVersion,
	)
	if err != nil {
		return fmt.Errorf("insert estimation record: %w", err)
	}
	return nil
}

// GetEstimation fetches a row by id, unchanged since it was saved.
func (s *Store) GetEstimation(id string) (record.EstimationRecord, bool, error) {
	row := s.db.QueryRow(`
		SELECT id, work_item_id, timestamp, specification_version, planning_version,
			repository_revision, features_json, prediction_json, risks_json,
			assumptions_json, model_version, calibration_version
		FROM estimation_records WHERE id = ?`, id)

	var (
		rec                                                      record.EstimationRecord
		timestampStr                                             string
		specVersion, planVersion, repoRev                        sql.NullString
		featuresJSON, predictionJSON, risksJSON, assumptionsJSON string
	)

	err := row.Scan(&rec.ID, &rec.WorkItemID, &timestampStr, &specVersion, &planVersion,
		&repoRev, &featuresJSON, &predictionJSON, &risksJSON, &assumptionsJSON,
		&rec.ModelVersion, &rec.CalibrationVersion)
	if err == sql.ErrNoRows {
		return record.EstimationRecord{}, false, nil
	}
	if err != nil {
		return record.EstimationRecord{}, false, fmt.Errorf("query estimation record: %w", err)
	}

	rec.Timestamp, err = time.Parse(time.RFC3339Nano, timestampStr)
	if err != nil {
		return record.EstimationRecord{}, false, fmt.Errorf("parse timestamp: %w", err)
	}
	if specVersion.Valid {
		v := specVersion.String
		rec.SpecificationVersion = &v
	}
	if planVersion.Valid {
		v := planVersion.String
		rec.PlanningVersion = &v
	}
	if repoRev.Valid {
		v := repoRev.String
		rec.RepositoryRevision = &v
	}
	if err := json.Unmarshal([]byte(featuresJSON), &rec.Features); err != nil {
		return record.EstimationRecord{}, false, fmt.Errorf("unmarshal features: %w", err)
	}
	if err := json.Unmarshal([]byte(predictionJSON), &rec.Prediction); err != nil {
		return record.EstimationRecord{}, false, fmt.Errorf("unmarshal prediction: %w", err)
	}
	if err := json.Unmarshal([]byte(risksJSON), &rec.Risks); err != nil {
		return record.EstimationRecord{}, false, fmt.Errorf("unmarshal risks: %w", err)
	}
	if err := json.Unmarshal([]byte(assumptionsJSON), &rec.Assumptions); err != nil {
		return record.EstimationRecord{}, false, fmt.Errorf("unmarshal assumptions: %w", err)
	}

	return rec, true, nil
}

// SaveOutcome inserts a new outcome_records row. The UNIQUE constraint on
// estimation_id (migrations.go) rejects a second outcome for the same
// estimation at the database level.
func (s *Store) SaveOutcome(o record.OutcomeRecord) error {
	_, err := s.db.Exec(`
		INSERT INTO outcome_records (
			id, estimation_id, actual_human_effort, actual_ai_usage, actual_ai_cost_usd,
			actual_lead_time_days, actual_verification_effort, actual_integration_effort,
			rework, incidents, completion_timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.ID, o.EstimationID, o.ActualHumanEffort, o.ActualAIUsage, o.ActualAICostUSD,
		o.ActualLeadTimeDays, o.ActualVerificationEffort, o.ActualIntegrationEffort,
		o.Rework, o.Incidents, o.CompletionTimestamp.Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("insert outcome record: %w", err)
	}
	return nil
}

// GetOutcome fetches the outcome linked to estimationID, if any.
func (s *Store) GetOutcome(estimationID string) (record.OutcomeRecord, bool, error) {
	row := s.db.QueryRow(`
		SELECT id, estimation_id, actual_human_effort, actual_ai_usage, actual_ai_cost_usd,
			actual_lead_time_days, actual_verification_effort, actual_integration_effort,
			rework, incidents, completion_timestamp
		FROM outcome_records WHERE estimation_id = ?`, estimationID)

	var o record.OutcomeRecord
	var completionStr string
	err := row.Scan(&o.ID, &o.EstimationID, &o.ActualHumanEffort, &o.ActualAIUsage,
		&o.ActualAICostUSD, &o.ActualLeadTimeDays, &o.ActualVerificationEffort,
		&o.ActualIntegrationEffort, &o.Rework, &o.Incidents, &completionStr)
	if err == sql.ErrNoRows {
		return record.OutcomeRecord{}, false, nil
	}
	if err != nil {
		return record.OutcomeRecord{}, false, fmt.Errorf("query outcome record: %w", err)
	}

	o.CompletionTimestamp, err = time.Parse(time.RFC3339Nano, completionStr)
	if err != nil {
		return record.OutcomeRecord{}, false, fmt.Errorf("parse completion timestamp: %w", err)
	}
	return o, true, nil
}
