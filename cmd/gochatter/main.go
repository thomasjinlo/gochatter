package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"github.com/thomasjinlo/gochatter/internal/auth"
	"github.com/thomasjinlo/gochatter/internal/client"
	"github.com/thomasjinlo/gochatter/internal/network"
)

func main() {
    currDir, _ := os.Getwd()
    viper.SetConfigName("config")
    viper.AddConfigPath(currDir + "/")
    err := viper.ReadInConfig()
    if err != nil { // Handle errors reading the config file
        log.Fatal("fatal error config file: %w", err)
    }

    serverConf := viper.Sub("server")

    switch os.Args[1] {
    case "server":
        tokenVerifier := auth.TokenVerifierFunc(auth.VerifyWithJWKS)
        err := http.ListenAndServeTLS(
            serverConf.GetString("port"),
            serverConf.GetString("certFile"),
            serverConf.GetString("keyFile"),
            http.HandlerFunc(network.HandleNewConnection(tokenVerifier)))
        if err != nil {
            log.Fatal("HTTP Server error:", err)
        }
    case "client":
        tlsConfig := &tls.Config{
            CipherSuites: []uint16{
                tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            },
        }
        dialer := &websocket.Dialer{TLSClientConfig: tlsConfig}

        authConfig := viper.Sub("auth")
        retrieveWithClientSecret := auth.RetrieveWithClientSecret(authConfig)
        tokenRetriever := auth.TokenRetrieverFunc(retrieveWithClientSecret)
        client := client.NewClient(
            serverConf.GetString("domain_name"),
            dialer,
            tokenRetriever,
        )
        client.Connect()
    }
}
