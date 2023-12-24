package util

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

// ComputeSHA256HMAC computes an HMAC for a given msg using SHA-256 hashing
// algorithm and returns corresponding checksum.
func ComputeSHA256HMAC(msg, key []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(msg)
	return hash.Sum(nil)
}

// NewKeyHMAC generates a cryptographic key with the specified byte size.
// The key can be used to compute an HMAC.
func NewKeyHMAC(size int) ([]byte, error) {
	if size < 1 {
		return nil, errors.New("size must be >= 1")
	}

	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}
