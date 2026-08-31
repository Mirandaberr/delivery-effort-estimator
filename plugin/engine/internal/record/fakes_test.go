package record

// In-memory fakes for unit-testing Service without touching real storage.
// The SQLite-backed implementation is exercised separately in
// internal/storage/sqlite.

type fakeEstimationRepo struct {
	byID map[string]EstimationRecord
}

func newFakeEstimationRepo() *fakeEstimationRepo {
	return &fakeEstimationRepo{byID: map[string]EstimationRecord{}}
}

func (f *fakeEstimationRepo) SaveEstimation(r EstimationRecord) error {
	f.byID[r.ID] = r
	return nil
}

func (f *fakeEstimationRepo) GetEstimation(id string) (EstimationRecord, bool, error) {
	r, ok := f.byID[id]
	return r, ok, nil
}

type fakeOutcomeRepo struct {
	byEstimationID map[string]OutcomeRecord
}

func newFakeOutcomeRepo() *fakeOutcomeRepo {
	return &fakeOutcomeRepo{byEstimationID: map[string]OutcomeRecord{}}
}

func (f *fakeOutcomeRepo) SaveOutcome(o OutcomeRecord) error {
	f.byEstimationID[o.EstimationID] = o
	return nil
}

func (f *fakeOutcomeRepo) GetOutcome(estimationID string) (OutcomeRecord, bool, error) {
	o, ok := f.byEstimationID[estimationID]
	return o, ok, nil
}
