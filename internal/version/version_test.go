package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalCommit := Commit

	// Test with commit
	Version = "v1.0.0"
	Commit = "abc123"
	result := String()
	expected := "v1.0.0 (abc123)"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	// Test without commit
	Commit = ""
	result = String()
	expected = "v1.0.0"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	// Test with default dev version
	Version = "dev"
	result = String()
	expected = "dev"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	// Restore original values
	Version = originalVersion
	Commit = originalCommit
}

func TestDefaultValues(t *testing.T) {
	// Test that default values are reasonable
	if Version == "" {
		t.Error("Version should not be empty")
	}

	result := String()
	if result == "" {
		t.Error("String() should not return empty string")
	}

	// Should contain the version
	if !strings.Contains(result, Version) {
		t.Errorf("String() result should contain version, got %s", result)
	}
}
