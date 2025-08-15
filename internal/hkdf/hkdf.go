package hkdf

import (
	"crypto/hmac"
	"crypto/sha256"
)

// Extract returns a pseudorandom key.
func Extract(salt, ikm []byte) []byte {
	h := hmac.New(sha256.New, salt)
	h.Write(ikm)
	return h.Sum(nil)
}

// Expand derives output keying material of length L.
func Expand(prk, info []byte, length int) []byte {
	if length <= 0 {
		return nil
	}
	var (
		t   []byte
		out []byte
		ctr byte = 1
	)
	for len(out) < length {
		h := hmac.New(sha256.New, prk)
		h.Write(t)
		h.Write(info)
		h.Write([]byte{ctr})
		t = h.Sum(nil)
		need := length - len(out)
		if need > len(t) {
			need = len(t)
		}
		out = append(out, t[:need]...)
		ctr++
	}
	return out
}

// Derive convenience (Extract then Expand).
func Derive(salt, ikm, info []byte, length int) []byte {
	return Expand(Extract(salt, ikm), info, length)
}
