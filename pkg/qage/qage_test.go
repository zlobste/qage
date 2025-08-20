package qage

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"
	"testing"

	"filippo.io/age"
)

func TestNewIdentity(t *testing.T) {
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	if id.Suite() != HybridX25519MLKEM768 {
		t.Errorf("expected suite %v, got %v", HybridX25519MLKEM768, id.Suite())
	}

	// Test that recipient is cached
	r1 := id.Recipient()
	r2 := id.Recipient()
	if r1 != r2 {
		t.Error("expected recipient to be cached")
	}

	if r1.Suite() != HybridX25519MLKEM768 {
		t.Errorf("expected recipient suite %v, got %v", HybridX25519MLKEM768, r1.Suite())
	}
}

func TestNewIdentityWithConfig(t *testing.T) {
	cfg := Config{Suite: HybridX25519MLKEM768}
	id, err := NewIdentityWithConfig(cfg)
	if err != nil {
		t.Fatalf("NewIdentityWithConfig failed: %v", err)
	}

	if id.Suite() != HybridX25519MLKEM768 {
		t.Errorf("expected suite %v, got %v", HybridX25519MLKEM768, id.Suite())
	}

	// Test with empty config (should use default)
	cfg = Config{}
	id, err = NewIdentityWithConfig(cfg)
	if err != nil {
		t.Fatalf("NewIdentityWithConfig with empty config failed: %v", err)
	}

	if id.Suite() != HybridX25519MLKEM768 {
		t.Errorf("expected default suite %v, got %v", HybridX25519MLKEM768, id.Suite())
	}

	// Test with unsupported suite
	cfg = Config{Suite: Suite(99)}
	_, err = NewIdentityWithConfig(cfg)
	if err == nil {
		t.Error("expected error for unsupported suite")
	}
}

func TestIdentityStringEncoding(t *testing.T) {
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	// Test String() method
	idStr, err := id.String()
	if err != nil {
		t.Fatalf("String() failed: %v", err)
	}

	if !strings.HasPrefix(idStr, "qagseck1") {
		t.Errorf("expected identity string to start with 'qagseck1', got %s", idStr)
	}

	// Test parsing back
	parsed, err := ParseIdentity(idStr)
	if err != nil {
		t.Fatalf("ParseIdentity failed: %v", err)
	}

	if parsed.Suite() != id.Suite() {
		t.Errorf("suite mismatch after parsing")
	}

	if parsed.x25519Secret != id.x25519Secret {
		t.Errorf("X25519 secret mismatch after parsing")
	}

	if !bytes.Equal(parsed.mlkemSecret, id.mlkemSecret) {
		t.Errorf("ML-KEM secret mismatch after parsing")
	}
}

func TestIdentityFileFormat(t *testing.T) {
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	// Test FormatFile without comment
	formatted, err := id.FormatFile("")
	if err != nil {
		t.Fatalf("FormatFile failed: %v", err)
	}

	if !strings.HasPrefix(formatted, "QAGE-SECRET-KEY-1 ") {
		t.Errorf("expected formatted string to start with 'QAGE-SECRET-KEY-1 ', got %s", formatted)
	}

	// Test FormatFile with comment
	comment := "test key"
	formatted, err = id.FormatFile(comment)
	if err != nil {
		t.Fatalf("FormatFile with comment failed: %v", err)
	}

	if !strings.Contains(formatted, comment) {
		t.Errorf("expected formatted string to contain comment '%s', got %s", comment, formatted)
	}

	// Test parsing back
	parsed, parsedComment, err := ParseIdentityFile(formatted)
	if err != nil {
		t.Fatalf("ParseIdentityFile failed: %v", err)
	}

	if parsedComment != comment {
		t.Errorf("expected comment '%s', got '%s'", comment, parsedComment)
	}

	if parsed.Suite() != id.Suite() {
		t.Errorf("suite mismatch after parsing")
	}
}

func TestRecipientStringEncoding(t *testing.T) {
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	r := id.Recipient()
	
	// Test String() method
	rStr, err := r.String()
	if err != nil {
		t.Fatalf("Recipient String() failed: %v", err)
	}

	if !strings.HasPrefix(rStr, "qage1") {
		t.Errorf("expected recipient string to start with 'qage1', got %s", rStr)
	}

	// Test parsing back
	parsed, err := ParseRecipient(rStr)
	if err != nil {
		t.Fatalf("ParseRecipient failed: %v", err)
	}

	if parsed.Suite() != r.Suite() {
		t.Errorf("suite mismatch after parsing")
	}

	if parsed.x25519Pub != r.x25519Pub {
		t.Errorf("X25519 public mismatch after parsing")
	}

	if !bytes.Equal(parsed.mlkemPub, r.mlkemPub) {
		t.Errorf("ML-KEM public mismatch after parsing")
	}
}

func TestAgeIntegration(t *testing.T) {
	// Create test data
	plaintext := make([]byte, 1024)
	if _, err := rand.Read(plaintext); err != nil {
		t.Fatalf("failed to generate test data: %v", err)
	}

	// Generate identity and recipient
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	r := id.Recipient()

	// Test encryption
	var encrypted bytes.Buffer
	w, err := age.Encrypt(&encrypted, r)
	if err != nil {
		t.Fatalf("age.Encrypt failed: %v", err)
	}

	if _, err := w.Write(plaintext); err != nil {
		t.Fatalf("writing plaintext failed: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("closing age writer failed: %v", err)
	}

	// Test decryption
	decReader, err := age.Decrypt(&encrypted, id)
	if err != nil {
		t.Fatalf("age.Decrypt failed: %v", err)
	}

	decrypted, err := io.ReadAll(decReader)
	if err != nil {
		t.Fatalf("reading decrypted data failed: %v", err)
	}

	// Verify
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted data doesn't match original")
	}
}

func TestWrapUnwrap(t *testing.T) {
	// Generate identity and recipient
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	r := id.Recipient()

	// Generate test file key
	fileKey := make([]byte, 32)
	if _, err := rand.Read(fileKey); err != nil {
		t.Fatalf("failed to generate file key: %v", err)
	}

	// Test wrap
	stanzas, err := r.Wrap(fileKey)
	if err != nil {
		t.Fatalf("Wrap failed: %v", err)
	}

	if len(stanzas) != 1 {
		t.Fatalf("expected 1 stanza, got %d", len(stanzas))
	}

	stanza := stanzas[0]
	if stanza.Type != "qage" {
		t.Errorf("expected stanza type 'qage', got '%s'", stanza.Type)
	}

	if len(stanza.Args) != 1 || stanza.Args[0] != "h1" {
		t.Errorf("expected args ['h1'], got %v", stanza.Args)
	}

	// Test unwrap
	unwrapped, err := id.Unwrap(stanzas)
	if err != nil {
		t.Fatalf("Unwrap failed: %v", err)
	}

	if !bytes.Equal(fileKey, unwrapped) {
		t.Errorf("unwrapped key doesn't match original")
	}

	// Test UnwrapStanza directly
	unwrapped2, err := id.UnwrapStanza(stanza)
	if err != nil {
		t.Fatalf("UnwrapStanza failed: %v", err)
	}

	if !bytes.Equal(fileKey, unwrapped2) {
		t.Errorf("UnwrapStanza result doesn't match original")
	}
}

func TestSuiteString(t *testing.T) {
	suite := HybridX25519MLKEM768
	expected := "X25519+ML-KEM-768"
	if suite.String() != expected {
		t.Errorf("expected %s, got %s", expected, suite.String())
	}

	// Test unknown suite
	unknown := Suite(99)
	result := unknown.String()
	if !strings.Contains(result, "Suite(99)") {
		t.Errorf("expected unknown suite string to contain 'Suite(99)', got %s", result)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Suite != HybridX25519MLKEM768 {
		t.Errorf("expected default suite %v, got %v", HybridX25519MLKEM768, cfg.Suite)
	}
}

func TestInvalidStanza(t *testing.T) {
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	// Test with invalid stanza type
	invalidStanza := &age.Stanza{
		Type: "invalid",
		Args: []string{"h1"},
		Body: []byte("invalid"),
	}

	_, err = id.UnwrapStanza(invalidStanza)
	if err == nil {
		t.Error("expected error for invalid stanza")
	}

	// Test with empty stanzas
	_, err = id.Unwrap([]*age.Stanza{})
	if err != age.ErrIncorrectIdentity {
		t.Errorf("expected ErrIncorrectIdentity, got %v", err)
	}
}

func TestUnsupportedSuite(t *testing.T) {
	// Test unsupported suite in wrap
	r := &Recipient{
		suite: Suite(99),
	}

	_, err := r.Wrap([]byte("test"))
	if err == nil {
		t.Error("expected error for unsupported suite in wrap")
	}

	// Test unsupported suite in unwrap
	id := &Identity{
		suite: Suite(99),
	}

	stanza := &age.Stanza{
		Type: "qage",
		Args: []string{"h1"},
		Body: []byte("test"),
	}

	_, err = id.UnwrapStanza(stanza)
	if err == nil {
		t.Error("expected error for unsupported suite in unwrap")
	}
}

func TestRecipientCaching(t *testing.T) {
	id, err := NewIdentity()
	if err != nil {
		t.Fatalf("NewIdentity failed: %v", err)
	}

	// Get recipient multiple times
	r1 := id.Recipient()
	r2 := id.Recipient()
	r3 := id.Recipient()

	// Should be the same instance
	if r1 != r2 || r2 != r3 {
		t.Error("recipient should be cached and return same instance")
	}

	// Test fallback case where cachedRecipient is nil
	id.cachedRecipient = nil
	r4 := id.Recipient()
	
	// Should return a zero recipient
	if r4.suite != id.suite {
		t.Error("fallback recipient should have same suite")
	}
	
	if r4.x25519Pub != [32]byte{} {
		t.Error("fallback recipient should have zero X25519 public key")
	}
	
	if r4.mlkemPub != nil {
		t.Error("fallback recipient should have nil ML-KEM public key")
	}
}
