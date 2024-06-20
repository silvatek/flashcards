package test

import (
	"testing"
)

func TestBasicScanning(t *testing.T) {
	source := "The quick `brown fox` jumped over the lazy dog"

	step := ParseStep(source)

	if step.Comparable != "The quick * jumped over the lazy dog" {
		t.Errorf("Unexpected comparable: %s", step.Comparable)
	}
	if len(step.Values) != 1 {
		t.Errorf("Unexpected value count: %d", len(step.Values))
	}
	if step.Values[0] != "brown fox" {
		t.Errorf("Unexpected value: %s", step.Values[0])
	}
}

func TestComparison(t *testing.T) {
	step1 := ParseStep("The `quick brown fox` jumped over the `lazy dog`")
	step2 := ParseStep("The `animal` jumped over the `barrier`")

	if !step1.Matches(step2) {
		t.Error("Parsed steps did not match")
	}
}
