package scenarios

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed catalog.json
var catalogJSON []byte

type Step struct {
	Service string `json:"service"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	Tone    string `json:"tone"`
}

type Scenario struct {
	ID              string `json:"id"`
	Index           string `json:"index"`
	Title           string `json:"title"`
	Kind            string `json:"kind"`
	FailureSlug     string `json:"failureSlug"`
	Summary         string `json:"summary"`
	InjectedFault   string `json:"injectedFault"`
	FundInvariant   string `json:"fundInvariant"`
	FirstAction     string `json:"firstAction"`
	RecoveryBasis   string `json:"recoveryBasis"`
	CurrentBoundary string `json:"currentBoundary"`
	GoTest          string `json:"goTest"`
	Steps           []Step `json:"steps"`
}

type Catalog struct {
	Version   int        `json:"version"`
	UpdatedAt string     `json:"updatedAt"`
	Scenarios []Scenario `json:"scenarios"`
}

func Load() (Catalog, error) {
	var catalog Catalog
	if err := json.Unmarshal(catalogJSON, &catalog); err != nil {
		return Catalog{}, fmt.Errorf("parse scenario catalog: %w", err)
	}
	return catalog, nil
}
