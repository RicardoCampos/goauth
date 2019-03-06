package oauth2

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	jwt "github.com/dgrijalva/jwt-go"
)

func getFutureInEpochTime() int64 {
	return time.Now().Add(time.Hour).UnixNano() / 1000000000
}

func createToken(scope string, expiry int64, tokenID string) string {
	nbfEpochTimeMs := GetNowInEpochTime()

	claims := TokenPayload{
		Scope: scope,
		StandardClaims: jwt.StandardClaims{
			Audience:  "youlot", // TODO: Really, this is optional, but will come from the client struct
			ExpiresAt: expiry,
			Id:        tokenID,
			IssuedAt:  nbfEpochTimeMs,
			Issuer:    "goauth.xom",   // TODO: This is a constant and should be the URI of the issuing server.
			NotBefore: nbfEpochTimeMs, // this should be the current time
			Subject:   "bar",          // Slightly off-piste here. We are not an authentication server so using the sub claim to identify the client.
		},
	}
	token, _ := SignToken(claims, LoadTestPrivateRsaKey("test_key"))
	return token
}
func TestValidateWithExpectedScope(t *testing.T) {
	expiry := getFutureInEpochTime()
	jti := uuid.New().String()
	token := createToken("read", expiry, jti)

	r, _ := NewReferenceToken(jti, "bar", expiry, token)
	publicKey := LoadTestPublicRsaKey("test_key.pub")
	result, err := ValidateToken(r, []string{"read"}, publicKey)
	assert.Nil(t, err)
	assert.Equal(t, result, true, "This should be valid")
}

func TestValidateWithInsufficientScope(t *testing.T) {
	expiry := getFutureInEpochTime()
	jti := uuid.New().String()
	token := createToken("read", expiry, jti)

	r, _ := NewReferenceToken(jti, "bar", expiry, token)
	publicKey := LoadTestPublicRsaKey("test_key.pub")
	result, err := ValidateToken(r, []string{"thisdoesnotexist"}, publicKey)
	assert.Equal(t, err, ErrInsufficientScope, "It should report insufficient scope")
	assert.Equal(t, result, false, "This should be invalid")
}
