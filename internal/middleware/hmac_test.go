package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/detod/best-wallet/internal/util"
)

func TestHMACVerifier_OK(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	sut := HMACVerifier(NewHMACKeyFetcherMock(map[string][]byte{keyID: key}))

	// Arrange: mock gin router.
	router := gin.New()
	method, path := "GET", "/sut"
	router.Handle(method, path, sut, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Arrange: signed request.
	body := []byte("some-body")
	msgToSign := []byte(fmt.Sprintf("%s%s", body, key))
	checksum := util.ComputeSHA256HMAC(msgToSign, key)
	signature := base64.StdEncoding.EncodeToString(checksum)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Add("BestWallet-Signature", signature)
	req.Header.Add("BestWallet-Key-ID", keyID)

	// Act: handle request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected response code %d, got %d", http.StatusNoContent, resp.Code)
	}
}

func TestHMACVerifier_MissingSignatureHeader(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	sut := HMACVerifier(NewHMACKeyFetcherMock(map[string][]byte{keyID: key}))

	// Arrange: mock gin router.
	router := gin.New()
	method, path := "GET", "/sut"
	router.Handle(method, path, sut, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Arrange: request without signature header, only keyID.
	req := httptest.NewRequest(method, path, bytes.NewBufferString("some-body"))
	req.Header.Add("BestWallet-Key-ID", keyID)

	// Act: handle request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected response code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestHMACVerifier_InvalidSignature(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	sut := HMACVerifier(NewHMACKeyFetcherMock(map[string][]byte{keyID: key}))

	// Arrange: mock gin router.
	router := gin.New()
	method, path := "GET", "/sut"
	router.Handle(method, path, sut, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Arrange: request with invalid signature.
	req := httptest.NewRequest(method, path, bytes.NewBufferString("some-body"))
	req.Header.Add("BestWallet-Signature", "invalid-signature")
	req.Header.Add("BestWallet-Key-ID", keyID)

	// Act: handle request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected response code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestHMACVerifier_MissingKeyIDHeader(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	sut := HMACVerifier(NewHMACKeyFetcherMock(map[string][]byte{keyID: key}))

	// Arrange: mock gin router.
	router := gin.New()
	method, path := "GET", "/sut"
	router.Handle(method, path, sut, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Arrange: request without keyID header, only signature.
	body := []byte("some-body")
	msgToSign := []byte(fmt.Sprintf("%s%s", body, key))
	checksum := util.ComputeSHA256HMAC(msgToSign, key)
	signature := base64.StdEncoding.EncodeToString(checksum)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Add("BestWallet-Signature", signature)

	// Act: handle request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected response code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestHMACVerifier_InvalidKeyID(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	sut := HMACVerifier(NewHMACKeyFetcherMock(map[string][]byte{keyID: key}))

	// Arrange: mock gin router.
	router := gin.New()
	method, path := "GET", "/sut"
	router.Handle(method, path, sut, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Arrange: request with invalid keyID.
	body := []byte("some-body")
	msgToSign := []byte(fmt.Sprintf("%s%s", body, key))
	checksum := util.ComputeSHA256HMAC(msgToSign, key)
	signature := base64.StdEncoding.EncodeToString(checksum)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Add("BestWallet-Signature", signature)
	req.Header.Add("BestWallet-Key-ID", "invalid-key")

	// Act: handle request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected response code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestHMACVerifier_EmptyBody(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	sut := HMACVerifier(NewHMACKeyFetcherMock(map[string][]byte{keyID: key}))

	// Arrange: mock gin router.
	router := gin.New()
	method, path := "GET", "/sut"
	router.Handle(method, path, sut, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Arrange: signed request with empty body.
	emptyBody := []byte{}
	msgToSign := []byte(fmt.Sprintf("%s%s", emptyBody, key))
	checksum := util.ComputeSHA256HMAC(msgToSign, key)
	signature := base64.StdEncoding.EncodeToString(checksum)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(emptyBody))
	req.Header.Add("BestWallet-Signature", signature)
	req.Header.Add("BestWallet-Key-ID", "invalid-key")

	// Act: handle request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected response code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func NewHMACKeyFetcherMock(m map[string][]byte) HMACKeyFetcher {
	return func(_ context.Context, keyID string) ([]byte, bool, error) {
		if key, ok := m[keyID]; ok {
			return key, true, nil
		}
		return nil, false, nil
	}
}
