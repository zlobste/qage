package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateX25519(t *testing.T) {
	priv1, pub1, err := GenerateX25519()
	if err != nil {
		t.Fatalf("GenerateX25519 failed: %v", err)
	}

	priv2, pub2, err := GenerateX25519()
	if err != nil {
		t.Fatalf("GenerateX25519 failed on second call: %v", err)
	}

	// Keys should be different
	if priv1 == priv2 {
		t.Error("private keys should be different")
	}

	if pub1 == pub2 {
		t.Error("public keys should be different")
	}

	// Keys should be 32 bytes
	if len(priv1) != 32 {
		t.Errorf("expected private key length 32, got %d", len(priv1))
	}

	if len(pub1) != 32 {
		t.Errorf("expected public key length 32, got %d", len(pub1))
	}
}

func TestGenerateMLKEM768(t *testing.T) {
	pub1, priv1, err := GenerateMLKEM768()
	if err != nil {
		t.Fatalf("GenerateMLKEM768 failed: %v", err)
	}

	pub2, priv2, err := GenerateMLKEM768()
	if err != nil {
		t.Fatalf("GenerateMLKEM768 failed on second call: %v", err)
	}

	// Keys should be different
	if bytes.Equal(priv1, priv2) {
		t.Error("private keys should be different")
	}

	if bytes.Equal(pub1, pub2) {
		t.Error("public keys should be different")
	}

	// Check expected lengths
	expectedPubLen := 1184  // ML-KEM-768 public key size
	expectedPrivLen := 2400 // ML-KEM-768 private key size

	if len(pub1) != expectedPubLen {
		t.Errorf("expected public key length %d, got %d", expectedPubLen, len(pub1))
	}

	if len(priv1) != expectedPrivLen {
		t.Errorf("expected private key length %d, got %d", expectedPrivLen, len(priv1))
	}
}
