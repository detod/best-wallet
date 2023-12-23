package utils

import (
	"encoding/hex"
	"testing"
)

func TestVerifySHA256HMAC_OK(t *testing.T) {
	// Arrange.
	msg := []byte("some-message")
	key, err := NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	checksum := ComputeSHA256HMAC(msg, key)

	// Act.
	ok := VerifySHA256HMAC(msg, key, checksum)

	// Assert.
	if !ok {
		t.Fatalf("expected true, got false")
	}
}

func TestVerifySHA256HMAC_DifferentMsg(t *testing.T) {
	// Arrange.
	key, err := NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	checksum := ComputeSHA256HMAC([]byte("some-message"), key)

	// Act.
	ok := VerifySHA256HMAC([]byte("some-other-message"), key, checksum)

	// Assert.
	if ok {
		t.Fatalf("expected false, got true")
	}
}

func TestComputeSHA256HMAC_Length(t *testing.T) {
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
		t.Fatalf("expected length to be %d, got %d", 32, len(checksum))
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

func TestNewKeyHMAC_ErrLengthLessThanOne(t *testing.T) {
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

func TestNewKeyHMAC_CustomLength(t *testing.T) {
	// Arrange.
	wantKeyLength := 64

	// Act.
	key, err := NewKeyHMAC(wantKeyLength)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Assert.
	if len(key) != wantKeyLength {
		t.Fatalf("expected length to be %d, got %d", wantKeyLength, len(key))
	}
}

func TestNewKeyHMAC_KeysDoNotRepeat(t *testing.T) {
	// Arrange.
	length := 10

	// Act.
	key1, err := NewKeyHMAC(length)
	if err != nil {
		t.Fatalf("error generating key1: %s", err)
	}
	key2, err := NewKeyHMAC(length)
	if err != nil {
		t.Fatalf("error generating key2: %s", err)
	}

	// Assert.
	if hex.EncodeToString(key1) == hex.EncodeToString(key2) {
		t.Fatalf("expected keys to be different but they're equal")
	}
}
