package hkdf

import (
	"encoding/hex"
	"testing"
)

// TestDeriveLength ensures HKDF derive returns correct length and determinism.
func TestDeriveLength(t *testing.T) {
	out1 := Derive([]byte("salt"), []byte("ikm"), []byte("info"), 48)
	out2 := Derive([]byte("salt"), []byte("ikm"), []byte("info"), 48)
	if len(out1) != 48 || len(out2) != 48 {
		t.Fatalf("unexpected length")
	}
	if hex.EncodeToString(out1) != hex.EncodeToString(out2) {
		t.Fatalf("non-deterministic")
	}
}

// TestExpandZero tests zero length returns nil.
func TestExpandZero(t *testing.T) {
	prk := Extract([]byte("s"), []byte("i"))
	if v := Expand(prk, []byte("info"), 0); v != nil {
		t.Fatalf("expected nil")
	}
}
