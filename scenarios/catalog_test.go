package scenarios

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

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

func TestCatalogMatchesPublishedV1Contract(t *testing.T) {
	catalog, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if catalog.Version != 1 {
		t.Fatalf("catalog version = %d, want 1", catalog.Version)
	}
	if _, err := time.Parse("2006-01-02", catalog.UpdatedAt); err != nil {
		t.Fatalf("updatedAt is not an ISO date: %v", err)
	}
	allowedTones := map[string]bool{"neutral": true, "active": true, "success": true, "warning": true, "danger": true, "recovered": true}
	for _, scenario := range catalog.Scenarios {
		if scenario.Kind != "baseline" && scenario.Kind != "failure" {
			t.Fatalf("scenario %s has invalid kind %q", scenario.ID, scenario.Kind)
		}
		for _, step := range scenario.Steps {
			if step.Service == "" || step.State == "" || step.Title == "" || step.Detail == "" || !allowedTones[step.Tone] {
				t.Fatalf("scenario %s has invalid step: %+v", scenario.ID, step)
			}
		}
	}
	schema, err := os.ReadFile("catalog.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var published map[string]any
	if err := json.Unmarshal(schema, &published); err != nil {
		t.Fatalf("catalog schema is invalid JSON: %v", err)
	}
	if published["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
		t.Fatalf("unexpected schema declaration: %v", published["$schema"])
	}
}
