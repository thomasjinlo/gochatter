package api

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/thomasjinlo/gochatter/internal/models"
)

func ListenAndServe() error {
	log.Print("[gochatter-api] Starting up API Server")
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[gochatter-api] root path %s", rootPath)
	rc := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	certPool := x509.NewCertPool()
	rootCA, err := os.ReadFile(os.Getenv("ROOTCA_CERT_PATH"))
	if err != nil {
		log.Printf("[gochatter-api] error reading ws cert: %v", err)
	}
	certPool.AppendCertsFromPEM(rootCA)
	tlsConfig := &tls.Config{RootCAs: certPool}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	hc := &http.Client{Transport: transport}
	pc := NewPushClient(rc, hc)
	ur := models.NewLocalUserRepository()
	mux := setupRoutes(pc, ur)

	return http.ListenAndServeTLS(
		":8443",
		os.Getenv("PUBLIC_CERT_PATH"),
		os.Getenv("PUBLIC_KEY_PATH"), mux)
}
