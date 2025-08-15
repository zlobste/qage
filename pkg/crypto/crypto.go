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

// MLKEM768Encapsulate performs ML-KEM-768 encapsulation.
func MLKEM768Encapsulate(publicKey []byte) (ciphertext []byte, sharedSecret []byte, err error) {
	var pk kyber768.PublicKey
	pk.Unpack(publicKey)

	ct := make([]byte, kyber768.CiphertextSize)
	ss := make([]byte, kyber768.SharedKeySize)
	pk.EncapsulateTo(ct, ss, nil)

	return ct, ss, nil
}

// MLKEM768Decapsulate performs ML-KEM-768 decapsulation.
func MLKEM768Decapsulate(privateKey []byte, ciphertext []byte) (sharedSecret []byte, err error) {
	var sk kyber768.PrivateKey
	sk.Unpack(privateKey)

	ss := make([]byte, kyber768.SharedKeySize)
	sk.DecapsulateTo(ss, ciphertext)

	return ss, nil
}

// X25519 performs X25519 ECDH.
func X25519(privateKey [32]byte, publicKey [32]byte) (sharedSecret [32]byte, err error) {
	curve := ecdh.X25519()
	priv, err := curve.NewPrivateKey(privateKey[:])
	if err != nil {
		return sharedSecret, err
	}

	pub, err := curve.NewPublicKey(publicKey[:])
	if err != nil {
		return sharedSecret, err
	}

	secret, err := priv.ECDH(pub)
	if err != nil {
		return sharedSecret, err
	}

	copy(sharedSecret[:], secret)
	return sharedSecret, nil
}
