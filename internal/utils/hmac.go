package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

// VerifySHA256HMAC returns true if checksum is a valid SHA-256 HMAC for msg.
func VerifySHA256HMAC(msg, key, checksum []byte) bool {
	return hmac.Equal(checksum, ComputeSHA256HMAC(msg, key))
}

// ComputeSHA256HMAC computes an HMAC for a given msg using SHA-256 hashing
// algorithm and returns corresponding checksum.
func ComputeSHA256HMAC(msg, key []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(msg)
	return hash.Sum(nil)
}

// NewKeyHMAC generates a secret cryptographic key for computing an HMAC.
func NewKeyHMAC(length int) ([]byte, error) {
	if length < 1 {
		return nil, errors.New("length must be >= 1")
	}

	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}
