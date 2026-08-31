// Command estimatorctl is the CLI contract described in
// specs/001-core-estimation-engine/contracts/cli.md.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jmirandev/delivery-effort-estimator/internal/estimation"
	"github.com/jmirandev/delivery-effort-estimator/internal/record"
	"github.com/jmirandev/delivery-effort-estimator/internal/storage/sqlite"
)

const defaultDBPath = "./data/estimator.db"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: estimatorctl <estimate|record-outcome|error-report|get> [flags]")
		os.Exit(1)
	}

	store, err := sqlite.Open(dbPath())
	if err != nil {
		fatal("storage_error", err)
	}
	defer store.Close()

	svc := record.NewService(store, store)

	cmd, args := os.Args[1], os.Args[2:]
	var runErr error
	switch cmd {
	case "estimate":
		runErr = runEstimate(svc, args)
	case "get":
		runErr = runGet(svc, args)
	case "record-outcome":
		runErr = runRecordOutcome(svc, args)
	case "error-report":
		runErr = runErrorReport(svc, args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}

	if runErr != nil {
		os.Exit(1)
	}
}

func dbPath() string {
	if v := os.Getenv("ESTIMATOR_DB_PATH"); v != "" {
		return v
	}
	return defaultDBPath
}

func runEstimate(svc *record.Service, args []string) error {
	fs := flag.NewFlagSet("estimate", flag.ExitOnError)
	workItem := fs.String("work-item", "", "work item id")
	featuresPath := fs.String("features", "", "path to features JSON, or '-' for stdin")
	_ = fs.Parse(args)

	if *workItem == "" {
		return emitError("missing_flag", map[string]string{"flag": "work-item"})
	}

	data, err := readInput(*featuresPath)
	if err != nil {
		return emitError("invalid_json", map[string]string{"detail": err.Error()})
	}

	rec, err := svc.EstimateFromJSON(*workItem, data)
	if err != nil {
		return handleEstimateError(err)
	}
	return printJSON(rec)
}

func runGet(svc *record.Service, args []string) error {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	estimationID := fs.String("estimation-id", "", "estimation id")
	_ = fs.Parse(args)

	if *estimationID == "" {
		return emitError("missing_flag", map[string]string{"flag": "estimation-id"})
	}

	rec, err := svc.GetEstimation(*estimationID)
	if err != nil {
		if errors.Is(err, record.ErrUnknownEstimation) {
			return emitError("unknown_estimation_id", map[string]string{"estimation_id": *estimationID})
		}
		return emitError("internal_error", map[string]string{"detail": err.Error()})
	}
	return printJSON(rec)
}

func runRecordOutcome(svc *record.Service, args []string) error {
	fs := flag.NewFlagSet("record-outcome", flag.ExitOnError)
	estimationID := fs.String("estimation-id", "", "estimation id")
	outcomePath := fs.String("outcome", "", "path to outcome JSON, or '-' for stdin")
	_ = fs.Parse(args)

	if *estimationID == "" {
		return emitError("missing_flag", map[string]string{"flag": "estimation-id"})
	}

	data, err := readInput(*outcomePath)
	if err != nil {
		return emitError("invalid_json", map[string]string{"detail": err.Error()})
	}

	outcome, err := svc.RecordOutcomeFromJSON(*estimationID, data)
	if err != nil {
		return handleOutcomeError(err, *estimationID)
	}
	return printJSON(outcome)
}

func runErrorReport(svc *record.Service, args []string) error {
	fs := flag.NewFlagSet("error-report", flag.ExitOnError)
	estimationID := fs.String("estimation-id", "", "estimation id")
	_ = fs.Parse(args)

	if *estimationID == "" {
		return emitError("missing_flag", map[string]string{"flag": "estimation-id"})
	}

	report, err := svc.ErrorReport(*estimationID)
	if err != nil {
		return handleReportError(err, *estimationID)
	}
	return printJSON(report)
}

func readInput(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func emitError(code string, extra map[string]string) error {
	payload := map[string]interface{}{"error": code}
	for k, v := range extra {
		payload[k] = v
	}
	b, _ := json.Marshal(payload)
	fmt.Fprintln(os.Stderr, string(b))
	return errors.New(code)
}

func handleEstimateError(err error) error {
	if field, ok := estimation.IsMissingFeatureError(err); ok {
		return emitError("missing_feature", map[string]string{"field": field})
	}
	if errors.Is(err, estimation.ErrInvalidJSON) {
		return emitError("invalid_json", map[string]string{"detail": err.Error()})
	}
	return emitError("internal_error", map[string]string{"detail": err.Error()})
}

func handleOutcomeError(err error, estimationID string) error {
	switch {
	case errors.Is(err, record.ErrUnknownEstimation):
		return emitError("unknown_estimation_id", map[string]string{"estimation_id": estimationID})
	case errors.Is(err, record.ErrOutcomeAlreadyExists):
		return emitError("outcome_already_recorded", map[string]string{"estimation_id": estimationID})
	case errors.Is(err, estimation.ErrInvalidJSON):
		return emitError("invalid_json", map[string]string{"detail": err.Error()})
	default:
		return emitError("internal_error", map[string]string{"detail": err.Error()})
	}
}

func handleReportError(err error, estimationID string) error {
	switch {
	case errors.Is(err, record.ErrUnknownEstimation):
		return emitError("unknown_estimation_id", map[string]string{"estimation_id": estimationID})
	case errors.Is(err, record.ErrNoOutcomeRecorded):
		return emitError("no_outcome_recorded", map[string]string{"estimation_id": estimationID})
	default:
		return emitError("internal_error", map[string]string{"detail": err.Error()})
	}
}

func fatal(code string, err error) {
	b, _ := json.Marshal(map[string]string{"error": code, "detail": err.Error()})
	fmt.Fprintln(os.Stderr, string(b))
	os.Exit(1)
}
