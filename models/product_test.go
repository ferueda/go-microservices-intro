package models

import (
	"testing"
)

func TestValidation(t *testing.T) {
	p := Product{Name: "adf", Price: 0.1, SKU: "abc-abc-abc"}

	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
}
