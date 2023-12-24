package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/detod/best-wallet/internal/util"
)

func TestHMACVerifier(t *testing.T) {
	// Arrange: secret key and id.
	keyID := "some-key-id"
	key, err := util.NewKeyHMAC(64)
	if err != nil {
		t.Fatalf("error generating a key: %s", err)
	}

	// Arrange: the SUT (system under test).
	mockKeyFetcher := HMACKeyFetcher(func(_ context.Context, id string) ([]byte, error) {
		if id != keyID {
			return nil, errors.New("key not found")
		}
		return key, nil
	})
	sut := HMACVerifier(mockKeyFetcher)

	// Arrange: mock gin router.
	router := gin.New()
	method := "GET"
	path := "/sut"
	handler := func(c *gin.Context) { c.Status(http.StatusNoContent) }
	router.Handle(method, path, sut, handler)

	// Arrange: a signed request.
	body := []byte("some-body")
	msgToSign := []byte(fmt.Sprintf("%s%s", body, key))
	checksum := util.ComputeSHA256HMAC(msgToSign, key)
	signature := base64.StdEncoding.EncodeToString(checksum)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Add("BestWallet-Signature", signature)
	req.Header.Add("BestWallet-Key-ID", keyID)

	// Act: handle signed request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert.
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected response code %d, got %d", http.StatusNoContent, resp.Code)
	}
}
