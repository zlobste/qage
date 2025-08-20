package encoding

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncodeDecodeRecipient(t *testing.T) {
	// Create test recipient
	r := &Recipient{
		Suite:     HybridX25519MLKEM768,
		X25519Pub: [32]byte{},
		MLKEMPub:  make([]byte, 1184), // ML-KEM-768 public key size
	}

	// Fill with random data
	if _, err := rand.Read(r.X25519Pub[:]); err != nil {
		t.Fatalf("failed to generate random X25519 key: %v", err)
	}
	if _, err := rand.Read(r.MLKEMPub); err != nil {
		t.Fatalf("failed to generate random ML-KEM key: %v", err)
	}

	// Encode
	encoded, err := EncodeRecipient(r)
	if err != nil {
		t.Fatalf("EncodeRecipient failed: %v", err)
	}

	// Should start with correct HRP
	if encoded[:5] != "qage1" {
		t.Errorf("expected encoded string to start with 'qage1', got %s", encoded[:5])
	}

	// Decode
	decoded, err := ParseRecipient(encoded)
	if err != nil {
		t.Fatalf("ParseRecipient failed: %v", err)
	}

	// Verify fields match
	if decoded.Suite != r.Suite {
		t.Errorf("suite mismatch: expected %v, got %v", r.Suite, decoded.Suite)
	}

	if decoded.X25519Pub != r.X25519Pub {
		t.Errorf("X25519 public key mismatch")
	}

	if !bytes.Equal(decoded.MLKEMPub, r.MLKEMPub) {
		t.Errorf("ML-KEM public key mismatch")
	}
}

func TestEncodeDecodeIdentity(t *testing.T) {
	// Create test identity
	id := &Identity{
		Suite:        HybridX25519MLKEM768,
		X25519Secret: [32]byte{},
		MLKEMSecret:  make([]byte, 2400), // ML-KEM-768 private key size
	}

	// Fill with random data
	if _, err := rand.Read(id.X25519Secret[:]); err != nil {
		t.Fatalf("failed to generate random X25519 key: %v", err)
	}
	if _, err := rand.Read(id.MLKEMSecret); err != nil {
		t.Fatalf("failed to generate random ML-KEM key: %v", err)
	}

	// Encode
	encoded, err := EncodeIdentity(id)
	if err != nil {
		t.Fatalf("EncodeIdentity failed: %v", err)
	}

	// Should start with correct HRP
	if encoded[:8] != "qagseck1" {
		t.Errorf("expected encoded string to start with 'qagseck1', got %s", encoded[:8])
	}

	// Decode
	decoded, err := ParseIdentity(encoded)
	if err != nil {
		t.Fatalf("ParseIdentity failed: %v", err)
	}

	// Verify fields match
	if decoded.Suite != id.Suite {
		t.Errorf("suite mismatch: expected %v, got %v", id.Suite, decoded.Suite)
	}

	if decoded.X25519Secret != id.X25519Secret {
		t.Errorf("X25519 secret key mismatch")
	}

	if !bytes.Equal(decoded.MLKEMSecret, id.MLKEMSecret) {
		t.Errorf("ML-KEM secret key mismatch")
	}
}

func TestFormatParseIdentityFile(t *testing.T) {
	// Create test identity
	id := &Identity{
		Suite:        HybridX25519MLKEM768,
		X25519Secret: [32]byte{1, 2, 3},
		MLKEMSecret:  make([]byte, 2400),
	}

	// Test without comment
	formatted, err := FormatIdentityFile(id, "")
	if err != nil {
		t.Fatalf("FormatIdentityFile without comment failed: %v", err)
	}

	parsedId, comment, err := ParseIdentityFile(formatted)
	if err != nil {
		t.Fatalf("ParseIdentityFile failed: %v", err)
	}

	if comment != "" {
		t.Errorf("expected empty comment, got '%s'", comment)
	}

	if parsedId.Suite != HybridX25519MLKEM768 {
		t.Errorf("expected suite %v, got %v", HybridX25519MLKEM768, parsedId.Suite)
	}

	// Test with comment
	formatted, err = FormatIdentityFile(id, "test comment")
	if err != nil {
		t.Fatalf("FormatIdentityFile with comment failed: %v", err)
	}

	parsedId2, comment2, err := ParseIdentityFile(formatted)
	if err != nil {
		t.Fatalf("ParseIdentityFile with comment failed: %v", err)
	}

	if comment2 != "test comment" {
		t.Errorf("expected comment 'test comment', got '%s'", comment2)
	}

	if parsedId2.Suite != HybridX25519MLKEM768 {
		t.Errorf("expected suite %v, got %v", HybridX25519MLKEM768, parsedId2.Suite)
	}
}

func TestParseRecipientErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid bech32", "invalid"},
		{"wrong HRP", "other1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8sf5e"},
		{"empty data", "qage1qqqqypcrgm"},
		{"unsupported suite", "qage1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvh6yhs"},
		{"wrong data length", "qage1sqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmk4m88"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRecipient(tt.input)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestParseIdentityErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid bech32", "invalid"},
		{"wrong HRP", "other1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8sf5e"},
		{"empty data", "qagseck1qqqqvjekrw"},
		{"unsupported suite", "qagseck1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqfq6hsx"},
		{"wrong data length", "qagseck1sqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqdvhf3s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseIdentity(tt.input)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestEncodeUnsupportedSuite(t *testing.T) {
	// Test recipient with unsupported suite
	r := &Recipient{
		Suite: Suite(99),
	}

	_, err := EncodeRecipient(r)
	if err == nil {
		t.Error("expected error for unsupported recipient suite")
	}

	// Test identity with unsupported suite
	id := &Identity{
		Suite: Suite(99),
	}

	_, err = EncodeIdentity(id)
	if err == nil {
		t.Error("expected error for unsupported identity suite")
	}
}

func TestParseIdentityFileErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty line", ""},
		{"comment line", "# this is a comment"},
		{"wrong prefix", "OTHER-SECRET-KEY-1 qagseck1..."},
		{"missing key", "QAGE-SECRET-KEY-1"},
		{"invalid key", "QAGE-SECRET-KEY-1 invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseIdentityFile(tt.input)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
