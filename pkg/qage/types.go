package qage

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"fmt"

	"filippo.io/age"
	kyber768 "github.com/cloudflare/circl/kem/kyber/kyber768"

	"github.com/zlobste/qage/internal/hkdf"
	"github.com/zlobste/qage/pkg/encoding"
)

// Identity represents a qage private identity for decryption.
type Identity struct {
	suite           Suite
	x25519Secret    [32]byte
	mlkemSecret     []byte
	cachedRecipient *Recipient
}

// Recipient represents a qage public recipient for encryption.
type Recipient struct {
	suite     Suite
	x25519Pub [32]byte
	mlkemPub  []byte
}

// Suite returns the cryptographic suite of the identity.
func (id *Identity) Suite() Suite {
	return id.suite
}

// Recipient returns the corresponding public recipient for this identity.
func (id *Identity) Recipient() *Recipient {
	if id.cachedRecipient != nil {
		return id.cachedRecipient
	}

	// This should never happen in normal usage since cachedRecipient is set during creation
	// Return a zero recipient as fallback
	return &Recipient{
		suite:     id.suite,
		x25519Pub: [32]byte{},
		mlkemPub:  nil,
	}
}

// String returns the bech32 encoding of the identity.
func (id *Identity) String() (string, error) {
	encId := &encoding.Identity{
		Suite:        encoding.Suite(id.suite),
		X25519Secret: id.x25519Secret,
		MLKEMSecret:  id.mlkemSecret,
	}
	return encoding.EncodeIdentity(encId)
}

// FormatFile returns the file format representation of the identity.
func (id *Identity) FormatFile(comment string) (string, error) {
	encId := &encoding.Identity{
		Suite:        encoding.Suite(id.suite),
		X25519Secret: id.x25519Secret,
		MLKEMSecret:  id.mlkemSecret,
	}
	return encoding.FormatIdentityFile(encId, comment)
}

// Suite returns the cryptographic suite of the recipient.
func (r *Recipient) Suite() Suite {
	return r.suite
}

// String returns the bech32 encoding of the recipient.
func (r *Recipient) String() (string, error) {
	encRec := &encoding.Recipient{
		Suite:     encoding.Suite(r.suite),
		X25519Pub: r.x25519Pub,
		MLKEMPub:  r.mlkemPub,
	}
	return encoding.EncodeRecipient(encRec)
}

// Age Integration Methods

// Ensure Recipient implements age.Recipient
var _ age.Recipient = (*Recipient)(nil)

// Wrap implements age.Recipient.
func (r *Recipient) Wrap(fileKey []byte) ([]*age.Stanza, error) {
	switch r.suite {
	case HybridX25519MLKEM768:
		return r.wrapHybridX25519MLKEM768(fileKey)
	default:
		return nil, fmt.Errorf("qage: unsupported suite %d", r.suite)
	}
}

func (r *Recipient) wrapHybridX25519MLKEM768(fileKey []byte) ([]*age.Stanza, error) {
	// Generate ephemeral X25519 key
	curve := ecdh.X25519()
	ephPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("qage: failed to generate ephemeral key: %w", err)
	}
	ephPub := ephPriv.PublicKey().Bytes()

	// ECDH with peer's X25519 public
	peerPub, err := curve.NewPublicKey(r.x25519Pub[:])
	if err != nil {
		return nil, fmt.Errorf("qage: invalid X25519 public key: %w", err)
	}
	z1, err := ephPriv.ECDH(peerPub)
	if err != nil {
		return nil, fmt.Errorf("qage: ECDH failed: %w", err)
	}

	// ML-KEM encapsulation
	var pk kyber768.PublicKey
	pk.Unpack(r.mlkemPub)
	ct := make([]byte, kyber768.CiphertextSize)
	z2 := make([]byte, kyber768.SharedKeySize)
	pk.EncapsulateTo(ct, z2, nil)

	// Hybrid KDF: derive wrap key from both shared secrets
	combined := make([]byte, 0, len(z1)+len(z2))
	combined = append(combined, z1...)
	combined = append(combined, z2...)
	wrapKey := hkdf.Derive(nil, combined, []byte("qage/wrap"), 32)

	// Encrypt file key using XOR
	encryptedKey := make([]byte, len(fileKey))
	for i := range fileKey {
		encryptedKey[i] = fileKey[i] ^ wrapKey[i%len(wrapKey)]
	}

	// Construct stanza body: ephPub || ct || encryptedKey
	body := make([]byte, 0, 32+kyber768.CiphertextSize+len(encryptedKey))
	body = append(body, ephPub...)
	body = append(body, ct...)
	body = append(body, encryptedKey...)

	stanza := &age.Stanza{
		Type: "qage",
		Args: []string{"h1"}, // hybrid version 1
		Body: body,
	}

	return []*age.Stanza{stanza}, nil
}

// Ensure Identity implements age.Identity
var _ age.Identity = (*Identity)(nil)

// Unwrap implements age.Identity.
func (id *Identity) Unwrap(stanzas []*age.Stanza) ([]byte, error) {
	for _, s := range stanzas {
		if s.Type == "qage" && len(s.Args) == 1 && s.Args[0] == "h1" {
			return id.unwrapStanza(s)
		}
	}
	return nil, age.ErrIncorrectIdentity
}

// UnwrapStanza unwraps a single stanza (used by tests and plugin).
func (id *Identity) UnwrapStanza(s *age.Stanza) ([]byte, error) {
	return id.unwrapStanza(s)
}

func (id *Identity) unwrapStanza(s *age.Stanza) ([]byte, error) {
	switch id.suite {
	case HybridX25519MLKEM768:
		return id.unwrapHybridX25519MLKEM768(s)
	default:
		return nil, fmt.Errorf("qage: unsupported suite %d", id.suite)
	}
}

func (id *Identity) unwrapHybridX25519MLKEM768(s *age.Stanza) ([]byte, error) {
	body := s.Body
	if len(body) < 32+kyber768.CiphertextSize {
		return nil, errors.New("qage: stanza too short")
	}

	// Parse stanza: ephPub || ct || encryptedKey
	ephPub := body[:32]
	ct := body[32 : 32+kyber768.CiphertextSize]
	encryptedKey := body[32+kyber768.CiphertextSize:]

	// ECDH with ephemeral public
	curve := ecdh.X25519()
	privKey, err := curve.NewPrivateKey(id.x25519Secret[:])
	if err != nil {
		return nil, fmt.Errorf("qage: invalid X25519 secret key: %w", err)
	}
	peerPub, err := curve.NewPublicKey(ephPub)
	if err != nil {
		return nil, fmt.Errorf("qage: invalid ephemeral public key: %w", err)
	}
	z1, err := privKey.ECDH(peerPub)
	if err != nil {
		return nil, fmt.Errorf("qage: ECDH failed: %w", err)
	}

	// ML-KEM decapsulation
	var sk kyber768.PrivateKey
	sk.Unpack(id.mlkemSecret)
	z2 := make([]byte, kyber768.SharedKeySize)
	sk.DecapsulateTo(z2, ct)

	// Hybrid KDF
	combined := make([]byte, 0, len(z1)+len(z2))
	combined = append(combined, z1...)
	combined = append(combined, z2...)
	wrapKey := hkdf.Derive(nil, combined, []byte("qage/wrap"), 32)

	// Decrypt file key using XOR
	fileKey := make([]byte, len(encryptedKey))
	for i := range encryptedKey {
		fileKey[i] = encryptedKey[i] ^ wrapKey[i%len(wrapKey)]
	}

	return fileKey, nil
}
