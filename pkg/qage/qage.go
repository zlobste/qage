// Package qage provides post-quantum hybrid encryption for age.
//
// This package implements a post-quantum secure key encapsulation mechanism (KEM)
// using a hybrid approach combining X25519 ECDH and ML-KEM-768. It provides
// drop-in recipients and identities compatible with filippo.io/age.
//
// # Basic Usage
//
//	// Generate a new identity
//	identity, err := qage.NewIdentity()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get the recipient for encryption
//	recipient := identity.Recipient()
//
//	// Use with age
//	out, _ := os.Create("secret.age")
//	w, _ := age.Encrypt(out, recipient)
//	w.Write([]byte("secret data"))
//	w.Close()
//
//	// Decrypt
//	in, _ := os.Open("secret.age")
//	r, _ := age.Decrypt(in, identity)
//	plaintext, _ := io.ReadAll(r)
//
// # Security
//
// This implementation provides post-quantum security through a hybrid KEM
// that combines:
//   - X25519 ECDH (classical security)
//   - ML-KEM-768 (post-quantum security)
//
// Both components must be broken to compromise the encryption.
package qage

import (
	"fmt"

	"github.com/zlobste/qage/pkg/crypto"
	"github.com/zlobste/qage/pkg/encoding"
)

// Suite identifies the cryptographic suite used.
type Suite uint8

const (
	// HybridX25519MLKEM768 combines X25519 ECDH with ML-KEM-768.
	HybridX25519MLKEM768 Suite = 1
)

// String returns the string representation of the suite.
func (s Suite) String() string {
	switch s {
	case HybridX25519MLKEM768:
		return "X25519+ML-KEM-768"
	default:
		return fmt.Sprintf("Suite(%d)", s)
	}
}

// Config specifies the cryptographic configuration.
type Config struct {
	Suite Suite
}

// DefaultConfig returns the default configuration using hybrid X25519+ML-KEM-768.
func DefaultConfig() Config {
	return Config{Suite: HybridX25519MLKEM768}
}

// NewIdentity generates a new identity with the default configuration.
func NewIdentity() (*Identity, error) {
	return NewIdentityWithConfig(DefaultConfig())
}

// NewIdentityWithConfig generates a new identity with the specified configuration.
func NewIdentityWithConfig(cfg Config) (*Identity, error) {
	if cfg.Suite == 0 {
		cfg = DefaultConfig()
	}

	switch cfg.Suite {
	case HybridX25519MLKEM768:
		return newHybridX25519MLKEM768Identity()
	default:
		return nil, fmt.Errorf("qage: unsupported suite %d", cfg.Suite)
	}
}

func newHybridX25519MLKEM768Identity() (*Identity, error) {
	// Generate X25519 keypair
	x25519Priv, x25519Pub, err := crypto.GenerateX25519()
	if err != nil {
		return nil, fmt.Errorf("qage: failed to generate X25519 key: %w", err)
	}

	// Generate ML-KEM-768 keypair
	mlkemPub, mlkemPriv, err := crypto.GenerateMLKEM768()
	if err != nil {
		return nil, fmt.Errorf("qage: failed to generate ML-KEM key: %w", err)
	}

	id := &Identity{
		suite:        HybridX25519MLKEM768,
		x25519Secret: x25519Priv,
		mlkemSecret:  mlkemPriv,
	}

	id.cachedRecipient = &Recipient{
		suite:     HybridX25519MLKEM768,
		x25519Pub: x25519Pub,
		mlkemPub:  mlkemPub,
	}

	return id, nil
}

// ParseRecipient parses a recipient string.
func ParseRecipient(recipientStr string) (*Recipient, error) {
	encRec, err := encoding.ParseRecipient(recipientStr)
	if err != nil {
		return nil, err
	}

	return &Recipient{
		suite:     Suite(encRec.Suite),
		x25519Pub: encRec.X25519Pub,
		mlkemPub:  encRec.MLKEMPub,
	}, nil
}

// ParseIdentity parses an identity from its bech32 representation.
func ParseIdentity(identityStr string) (*Identity, error) {
	encId, err := encoding.ParseIdentity(identityStr)
	if err != nil {
		return nil, err
	}

	return &Identity{
		suite:        Suite(encId.Suite),
		x25519Secret: encId.X25519Secret,
		mlkemSecret:  encId.MLKEMSecret,
	}, nil
}

// ParseIdentityFile parses an identity from a file line.
func ParseIdentityFile(line string) (*Identity, string, error) {
	encId, comment, err := encoding.ParseIdentityFile(line)
	if err != nil {
		return nil, "", err
	}

	id := &Identity{
		suite:        Suite(encId.Suite),
		x25519Secret: encId.X25519Secret,
		mlkemSecret:  encId.MLKEMSecret,
	}

	return id, comment, nil
}
