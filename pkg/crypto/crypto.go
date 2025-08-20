// Package crypto provides cryptographic primitives for qage.
package crypto

import (
	"crypto/ecdh"
	"crypto/rand"

	kyber768 "github.com/cloudflare/circl/kem/kyber/kyber768"
)

// GenerateX25519 generates a new X25519 keypair.
func GenerateX25519() (privateKey [32]byte, publicKey [32]byte, err error) {
	curve := ecdh.X25519()
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return privateKey, publicKey, err
	}

	pub := priv.PublicKey().Bytes()
	copy(privateKey[:], priv.Bytes())
	copy(publicKey[:], pub)

	return privateKey, publicKey, nil
}

// GenerateMLKEM768 generates a new ML-KEM-768 keypair.
func GenerateMLKEM768() (publicKey []byte, privateKey []byte, err error) {
	pk, sk, err := kyber768.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Pack keys into byte slices
	pkBytes := make([]byte, kyber768.PublicKeySize)
	skBytes := make([]byte, kyber768.PrivateKeySize)
	pk.Pack(pkBytes)
	sk.Pack(skBytes)

	return pkBytes, skBytes, nil
}
