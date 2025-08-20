package qage

import (
	"testing"
)

func TestSelftest(t *testing.T) {
	// Test that selftest passes
	if err := Selftest(); err != nil {
		t.Fatalf("Selftest failed: %v", err)
	}
}

func TestSelftestKeygenRoundtrip(t *testing.T) {
	// Test the individual selftest functions directly
	if err := SelftestKeygenRoundtrip(); err != nil {
		t.Fatalf("SelftestKeygenRoundtrip failed: %v", err)
	}
}

func TestSelftestAgeIntegration(t *testing.T) {
	if err := SelftestAgeIntegration(); err != nil {
		t.Fatalf("SelftestAgeIntegration failed: %v", err)
	}
}

func TestSelftestStringEncoding(t *testing.T) {
	if err := SelftestStringEncoding(); err != nil {
		t.Fatalf("SelftestStringEncoding failed: %v", err)
	}
}
