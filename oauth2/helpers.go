package oauth2

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetNowInEpochTime returns the epoch time normally used in JWT's
func GetNowInEpochTime() int64 {
	return time.Now().UnixNano() / 1000000000
}

func LoadTestPrivateRsaKey(keylocation string) *rsa.PrivateKey {
	k, err := ioutil.ReadFile(keylocation)
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(k)
	if err != nil {
		panic(err)
	}
	return key
}
func LoadTestPublicRsaKey(keylocation string) *rsa.PublicKey {
	k, err := ioutil.ReadFile(keylocation)
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(k)
	if err != nil {
		panic(err)
	}
	return key
}
