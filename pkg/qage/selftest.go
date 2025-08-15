package qage

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"filippo.io/age"
)

// Selftest runs internal validation tests.
func Selftest() error {
	// Test 1: Key generation and round-trip
	if err := selftestKeygenRoundtrip(); err != nil {
		return fmt.Errorf("keygen roundtrip: %w", err)
	}

	// Test 2: Age integration
	if err := selftestAgeIntegration(); err != nil {
		return fmt.Errorf("age integration: %w", err)
	}

	// Test 3: String encoding/decoding
	if err := selftestStringEncoding(); err != nil {
		return fmt.Errorf("string encoding: %w", err)
	}

	return nil
}

func selftestKeygenRoundtrip() error {
	// Generate identity
	id, err := NewIdentity()
	if err != nil {
		return err
	}

	// Encode and decode identity
	idStr, err := id.String()
	if err != nil {
		return err
	}

	id2, err := ParseIdentity(idStr)
	if err != nil {
		return err
	}

	// Verify they're equivalent
	if id.suite != id2.suite {
		return errors.New("suite mismatch")
	}
	if id.x25519Secret != id2.x25519Secret {
		return errors.New("X25519 secret mismatch")
	}
	if !bytes.Equal(id.mlkemSecret, id2.mlkemSecret) {
		return errors.New("ML-KEM secret mismatch")
	}

	// Test recipient
	r := id.Recipient()
	rStr, err := r.String()
	if err != nil {
		return err
	}

	r2, err := ParseRecipient(rStr)
	if err != nil {
		return err
	}

	if r.suite != r2.suite {
		return errors.New("recipient suite mismatch")
	}
	if r.x25519Pub != r2.x25519Pub {
		return errors.New("recipient X25519 public mismatch")
	}
	if !bytes.Equal(r.mlkemPub, r2.mlkemPub) {
		return errors.New("recipient ML-KEM public mismatch")
	}

	return nil
}

func selftestAgeIntegration() error {
	// Generate test data
	plaintext := make([]byte, 1024)
	if _, err := rand.Read(plaintext); err != nil {
		return err
	}

	// Generate identity and recipient
	id, err := NewIdentity()
	if err != nil {
		return err
	}
	r := id.Recipient()

	// Encrypt with age
	var encrypted bytes.Buffer
	w, err := age.Encrypt(&encrypted, r)
	if err != nil {
		return err
	}
	if _, writeErr := w.Write(plaintext); writeErr != nil {
		return writeErr
	}
	if closeErr := w.Close(); closeErr != nil {
		return closeErr
	}

	// Decrypt with age
	decrypted, err := age.Decrypt(&encrypted, id)
	if err != nil {
		return err
	}

	decryptedData, err := io.ReadAll(decrypted)
	if err != nil {
		return err
	}

	// Verify data
	if !bytes.Equal(plaintext, decryptedData) {
		return errors.New("plaintext mismatch")
	}

	return nil
}

func selftestStringEncoding() error {
	// Test file format
	id, err := NewIdentity()
	if err != nil {
		return err
	}

	fileStr, err := id.FormatFile("test comment")
	if err != nil {
		return err
	}

	id2, comment, err := ParseIdentityFile(fileStr)
	if err != nil {
		return err
	}

	if comment != "test comment" {
		return fmt.Errorf("comment mismatch: got %q, expected %q", comment, "test comment")
	}

	// Verify identity
	if id.suite != id2.suite {
		return errors.New("file format suite mismatch")
	}
	if id.x25519Secret != id2.x25519Secret {
		return errors.New("file format X25519 secret mismatch")
	}
	if !bytes.Equal(id.mlkemSecret, id2.mlkemSecret) {
		return errors.New("file format ML-KEM secret mismatch")
	}

	return nil
}
