package main

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ricardocampos/goauth/services"
)

func useDB() bool {
	use := os.Getenv("CONNECTION_STRING")
	if len(use) < 1 || use == "false" {
		return false
	}
	return true
}

func loadRsaKey() *rsa.PrivateKey {
	k, err := ioutil.ReadFile("sample_key")
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(k)
	if err != nil {
		panic(err)
	}
	return key
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	rsaKey := loadRsaKey()
	var svc services.OAuth2Service
	if !useDB() {
		svc = services.NewInMemoryOAuth2Service(logger, rsaKey)
	} else {
		svc = services.NewPostgresOAuth2Service(logger, os.Getenv("CONNECTION_STRING"), rsaKey)
	}

	// Wrap services in middleware for middlware goodness.
	svc = services.NewLoggingMiddleware(logger, svc)
	svc = services.NewInstrumentingMiddleware(svc)

	tokenHandler := services.NewTokenHandler(svc)
	validateHandler := services.NewValidateHandler(svc)

	http.Handle("/token", tokenHandler)
	http.Handle("/validate", validateHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
