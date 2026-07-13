package scenarios

import "testing"

func TestCatalogHasBaselineAndSixFailures(t *testing.T) {
	catalog, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Scenarios) != 7 {
		t.Fatalf("scenario count = %d, want 7", len(catalog.Scenarios))
	}
	ids := map[string]bool{}
	failures := 0
	for _, scenario := range catalog.Scenarios {
		if scenario.ID == "" || ids[scenario.ID] {
			t.Fatalf("invalid or duplicate scenario id %q", scenario.ID)
		}
		ids[scenario.ID] = true
		if len(scenario.Steps) < 4 || scenario.FundInvariant == "" || scenario.FirstAction == "" || scenario.RecoveryBasis == "" || scenario.GoTest == "" {
			t.Fatalf("scenario %s is missing recovery evidence", scenario.ID)
		}
		if scenario.Kind == "failure" {
			failures++
			if scenario.FailureSlug == "" {
				t.Fatalf("failure scenario %s has no site relation", scenario.ID)
			}
		}
	}
	if failures != 6 {
		t.Fatalf("failure count = %d, want 6", failures)
	}
}
