package encoding

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	buf := make([]byte, 64)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("failed to generate random bytes: %v", err)
	}
	s, err := Encode("qage", buf)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if s[:5] != "qage1" {
		t.Fatalf("prefix")
	}
	t.Logf("encoded=%s len=%d", s, len(s))
	hrp, out, err := Decode(s)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if hrp != "qage" {
		t.Fatalf("hrp")
	}
	if !bytes.Equal(buf, out) {
		t.Fatalf("roundtrip mismatch")
	}
}

func TestChecksumReject(t *testing.T) {
	s, _ := Encode("qage", []byte{0, 1, 2})
	// Corrupt last char
	bad := s[:len(s)-1] + "x"
	if _, _, err := Decode(bad); err == nil {
		t.Fatalf("expected error")
	}
}
