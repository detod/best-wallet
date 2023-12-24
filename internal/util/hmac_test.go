package util

import (
	"encoding/hex"
	"testing"
)

func TestComputeSHA256HMAC_Size(t *testing.T) {
	// Arrange.
	msg := []byte("some-message")
	key, err := NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Act.
	checksum := ComputeSHA256HMAC(msg, key)

	// Assert.
	if len(checksum) != 32 { // 256 bits = 32 bytes
		t.Fatalf("expected size to be %d, got %d", 32, len(checksum))
	}
}

func TestComputeSHA256HMAC_SameInputSameOutput(t *testing.T) {
	// Arrange.
	msg := []byte("some-message")
	key, err := NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Act.
	checksum := ComputeSHA256HMAC(msg, key)
	checksumRepeat := ComputeSHA256HMAC(msg, key)

	// Assert.
	if hex.EncodeToString(checksum) != hex.EncodeToString(checksumRepeat) {
		t.Fatalf("expected checksum to repeat for the same msg and key, but it did not")
	}
}

func TestNewKeyHMAC_ErrSizeLessThanOne(t *testing.T) {
	// Act.
	key, err := NewKeyHMAC(0)

	// Assert.
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if key != nil {
		t.Fatalf("expected key to be nil, got %s", key)
	}
}

func TestNewKeyHMAC_CorrectSize(t *testing.T) {
	// Arrange.
	wantSize := 64

	// Act.
	key, err := NewKeyHMAC(wantSize)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Assert.
	if len(key) != wantSize {
		t.Fatalf("expected size to be %d, got %d", wantSize, len(key))
	}
}

func TestNewKeyHMAC_KeysDoNotRepeat(t *testing.T) {
	// Arrange.
	size := 10

	// Act.
	key1, err := NewKeyHMAC(size)
	if err != nil {
		t.Fatalf("error generating key1: %s", err)
	}
	key2, err := NewKeyHMAC(size)
	if err != nil {
		t.Fatalf("error generating key2: %s", err)
	}

	// Assert.
	if hex.EncodeToString(key1) == hex.EncodeToString(key2) {
		t.Fatalf("expected keys to be different but they're equal")
	}
}
