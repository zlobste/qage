package encoding

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Key HRPs (Human Readable Parts)
const (
	HRPPublic = "qage"
	HRPSecret = "qagseck"
)

// Suite identifies the cryptographic suite.
type Suite uint8

const (
	HybridX25519MLKEM768 Suite = 1
)

// Recipient represents a qage recipient.
type Recipient struct {
	Suite     Suite
	X25519Pub [32]byte
	MLKEMPub  []byte
}

// Identity represents a qage identity.
type Identity struct {
	Suite        Suite
	X25519Secret [32]byte
	MLKEMSecret  []byte
}

// ParseRecipient parses a qage recipient from its bech32 encoding.
func ParseRecipient(recipientStr string) (*Recipient, error) {
	hrp, data, err := Decode(recipientStr)
	if err != nil {
		return nil, fmt.Errorf("qage: invalid recipient encoding: %w", err)
	}

	if hrp != HRPPublic {
		return nil, fmt.Errorf("qage: invalid recipient HRP %q, expected %q", hrp, HRPPublic)
	}

	if len(data) == 0 {
		return nil, errors.New("qage: empty recipient data")
	}

	// Parse version byte
	suite := Suite(data[0])
	data = data[1:]

	switch suite {
	case HybridX25519MLKEM768:
		return parseHybridX25519MLKEM768Recipient(data)
	default:
		return nil, fmt.Errorf("qage: unsupported suite %d", suite)
	}
}

func parseHybridX25519MLKEM768Recipient(data []byte) (*Recipient, error) {
	const expectedLen = 32 + 1184 // X25519 pub + ML-KEM-768 pub
	if len(data) != expectedLen {
		return nil, fmt.Errorf("qage: invalid hybrid recipient length %d, expected %d", len(data), expectedLen)
	}

	r := &Recipient{Suite: HybridX25519MLKEM768}
	copy(r.X25519Pub[:], data[:32])
	r.MLKEMPub = make([]byte, 1184)
	copy(r.MLKEMPub, data[32:])

	return r, nil
}

// ParseIdentity parses a qage identity from its bech32 encoding.
func ParseIdentity(identityStr string) (*Identity, error) {
	hrp, data, err := Decode(identityStr)
	if err != nil {
		return nil, fmt.Errorf("qage: invalid identity encoding: %w", err)
	}

	if hrp != HRPSecret {
		return nil, fmt.Errorf("qage: invalid identity HRP %q, expected %q", hrp, HRPSecret)
	}

	if len(data) == 0 {
		return nil, errors.New("qage: empty identity data")
	}

	// Parse version byte
	suite := Suite(data[0])
	data = data[1:]

	switch suite {
	case HybridX25519MLKEM768:
		return parseHybridX25519MLKEM768Identity(data)
	default:
		return nil, fmt.Errorf("qage: unsupported suite %d", suite)
	}
}

func parseHybridX25519MLKEM768Identity(data []byte) (*Identity, error) {
	const expectedLen = 32 + 2400 // X25519 priv + ML-KEM-768 priv
	if len(data) != expectedLen {
		return nil, fmt.Errorf("qage: invalid hybrid identity length %d, expected %d", len(data), expectedLen)
	}

	id := &Identity{Suite: HybridX25519MLKEM768}
	copy(id.X25519Secret[:], data[:32])
	id.MLKEMSecret = make([]byte, 2400)
	copy(id.MLKEMSecret, data[32:])

	return id, nil
}

// EncodeRecipient encodes a recipient to its bech32 representation.
func EncodeRecipient(r *Recipient) (string, error) {
	switch r.Suite {
	case HybridX25519MLKEM768:
		return encodeHybridX25519MLKEM768Recipient(r)
	default:
		return "", fmt.Errorf("qage: unsupported suite %d", r.Suite)
	}
}

func encodeHybridX25519MLKEM768Recipient(r *Recipient) (string, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(r.Suite))
	buf.Write(r.X25519Pub[:])
	buf.Write(r.MLKEMPub)

	return Encode(HRPPublic, buf.Bytes())
}

// EncodeIdentity encodes an identity to its bech32 representation.
func EncodeIdentity(id *Identity) (string, error) {
	switch id.Suite {
	case HybridX25519MLKEM768:
		return encodeHybridX25519MLKEM768Identity(id)
	default:
		return "", fmt.Errorf("qage: unsupported suite %d", id.Suite)
	}
}

func encodeHybridX25519MLKEM768Identity(id *Identity) (string, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(id.Suite))
	buf.Write(id.X25519Secret[:])
	buf.Write(id.MLKEMSecret)

	return Encode(HRPSecret, buf.Bytes())
}

// ParseIdentityFile parses an identity from the "QAGE-SECRET-KEY-1 <bech32> # comment" format.
func ParseIdentityFile(line string) (*Identity, string, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, "", errors.New("qage: empty or comment line")
	}

	if !strings.HasPrefix(line, "QAGE-SECRET-KEY-1 ") {
		return nil, "", errors.New("qage: invalid identity line format")
	}

	// Remove prefix
	line = strings.TrimPrefix(line, "QAGE-SECRET-KEY-1 ")

	// Split on first whitespace to separate key from comment
	parts := strings.SplitN(line, " ", 2)
	keyStr := parts[0]

	var comment string
	if len(parts) > 1 {
		comment = strings.TrimSpace(parts[1])
		if strings.HasPrefix(comment, "#") {
			comment = strings.TrimSpace(comment[1:])
		}
	}

	id, err := ParseIdentity(keyStr)
	if err != nil {
		return nil, "", err
	}

	return id, comment, nil
}

// FormatIdentityFile formats an identity for storage in a file.
func FormatIdentityFile(id *Identity, comment string) (string, error) {
	encoded, err := EncodeIdentity(id)
	if err != nil {
		return "", err
	}

	if comment != "" {
		return fmt.Sprintf("QAGE-SECRET-KEY-1 %s # %s", encoded, comment), nil
	}
	return fmt.Sprintf("QAGE-SECRET-KEY-1 %s", encoded), nil
}
