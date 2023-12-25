package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/detod/best-wallet/internal/util"
)

const sigHeader = "BestWallet-Signature"
const keyIDHeader = "BestWallet-Key-ID"

// HMACKeyFetcher can fetch an HMAC secret key by a corresponding key ID.
// The key can be used for computing a signature i.e. signing a request.
type HMACKeyFetcher func(ctx context.Context, keyID string) (key []byte, found bool, err error)

// HMACVerifier will check if the request was signed by an authorized client.
// Signing a request requires a secret key that is exchanged ahead of time
// between server and client.
//
// The signature is computed as follows:
// 1. message_to_sign = concat(request_body, secret_key)
// 2. checksum = hmac_sha256(message_to_sign, secret_key)
// 3. signature = base64_encode(checksum)
//
// The client then sends the signature as part of the request (in a header)
// thus proving its identity as well as the integrity of the request payload.
// This function will verify the signature by re-calculating the signature
// itself, using the same method as above, and compare it to the signature
// provided by the client. If the signatures are the same, verification is
// successful and the server can continue processing the request. If the
// signatures are not the same, verificaion failed and the request is denied.
//
// Note that this is only a demo implementation and it's missing a broader scope
// of the message to sign (URI, HTTP method, timestamp/nonce, headers).
func HMACVerifier(fetchKey HMACKeyFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		sig, keyID := c.GetHeader(sigHeader), c.GetHeader(keyIDHeader)
		if sig == "" || keyID == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		key, found, err := fetchKey(c.Request.Context(), keyID)
		if !found {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		msg := []byte(fmt.Sprintf("%s%s", body, key))
		sum := util.ComputeSHA256HMAC(msg, key)
		wantSig := base64.StdEncoding.EncodeToString(sum)

		if sig != wantSig {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
