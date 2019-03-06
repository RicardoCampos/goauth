package main

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ricardocampos/goauth/handlers"
	"github.com/ricardocampos/goauth/services"
)

func useDB() bool {
	use := os.Getenv("CONNECTION_STRING")
	if len(use) < 1 || use == "false" {
		return false
	}
	return true
}

func loadPrivateRsaKey() *rsa.PrivateKey {
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

func loadPublicRsaKey() *rsa.PublicKey {
	k, err := ioutil.ReadFile("sample_key.pub")
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(k)
	if err != nil {
		panic(err)
	}
	return key
}


func main() {
	r := mux.NewRouter()
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout
	rsaPrivateKey := loadPrivateRsaKey()
	rsaPublicKey := loadPublicRsaKey()
	var svc services.TokenService
	if !useDB() {
		svc = services.NewInMemoryTokenService(logger, rsaPrivateKey, rsaPublicKey)
	} else {
		svc = services.NewPostgresTokenService(logger, os.Getenv("CONNECTION_STRING"), rsaPrivateKey, rsaPublicKey)
	}
	tokenHandler := handlers.NewTokenHandler(logger, svc)
	validateHandler := handlers.NewValidateHandler(logger, svc)
	r.HandleFunc("/token", tokenHandler.HandleToken)
	r.HandleFunc("/validate", validateHandler.HandleValidation)
	r.Handle("/metrics", promhttp.Handler())
	http.Handle("/", r)
	logger.Info("msg", "HTTP", "addr", ":8080")
	logger.Info("err", http.ListenAndServe(":8080", nil))
}


