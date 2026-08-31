package record

import (
	"reflect"
	"strings"
	"testing"
)

// TestRepositoriesExposeNoMutationMethods guards Constitution Principle VIII
// (estimations/outcomes are append-only) at the interface level: it fails
// loudly if either repository interface ever grows an Update/Delete-shaped
// method.
func TestRepositoriesExposeNoMutationMethods(t *testing.T) {
	forbidden := []string{"Update", "Delete", "Remove", "Overwrite"}
	check := func(name string, ifaceType reflect.Type) {
		for i := 0; i < ifaceType.NumMethod(); i++ {
			m := ifaceType.Method(i)
			for _, bad := range forbidden {
				if strings.Contains(m.Name, bad) {
					t.Errorf("%s exposes mutation-shaped method %s", name, m.Name)
				}
			}
		}
	}
	check("EstimationRepository", reflect.TypeOf((*EstimationRepository)(nil)).Elem())
	check("OutcomeRepository", reflect.TypeOf((*OutcomeRepository)(nil)).Elem())
}
