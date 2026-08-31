package sqlite

import "database/sql"

// schema is append-only by construction: no ALTER/DROP of existing rows, and
// the application layer (internal/record) never issues UPDATE/DELETE against
// these tables (Constitution Principle VIII). The UNIQUE constraint on
// outcome_records.estimation_id enforces "at most one outcome per estimation"
// (FR-013) even if application logic has a bug.
const schema = `
CREATE TABLE IF NOT EXISTS estimation_records (
	id                       TEXT PRIMARY KEY,
	work_item_id             TEXT NOT NULL,
	timestamp                TEXT NOT NULL,
	specification_version    TEXT,
	planning_version         TEXT,
	repository_revision      TEXT,
	features_json            TEXT NOT NULL,
	prediction_json          TEXT NOT NULL,
	confidence               REAL NOT NULL,
	prediction_interval_p50  REAL NOT NULL,
	prediction_interval_p80  REAL NOT NULL,
	risks_json               TEXT NOT NULL,
	assumptions_json         TEXT NOT NULL,
	model_version            TEXT NOT NULL,
	calibration_version      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS outcome_records (
	id                          TEXT PRIMARY KEY,
	estimation_id               TEXT NOT NULL UNIQUE REFERENCES estimation_records(id),
	actual_human_effort         REAL NOT NULL,
	actual_ai_usage             REAL NOT NULL,
	actual_ai_cost_usd          REAL NOT NULL,
	actual_lead_time_days       REAL NOT NULL,
	actual_verification_effort  REAL NOT NULL,
	actual_integration_effort   REAL NOT NULL,
	rework                      REAL NOT NULL,
	incidents                   INTEGER NOT NULL,
	completion_timestamp        TEXT NOT NULL
);
`

func migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
