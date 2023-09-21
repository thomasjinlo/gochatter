package main

import (
	"crypto/tls"
	"crypto/x509"
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
    switch os.Args[1] {
    case "server":
        tokenVerifier := auth.TokenVerifierFunc(auth.VerifyWithJWKS)
        serverConf := viper.Sub("server")
        err := http.ListenAndServeTLS(
            serverConf.GetString("port"),
            serverConf.GetString("certFile"),
            serverConf.GetString("keyFile"),
            http.HandlerFunc(network.HandleNewConnection(tokenVerifier)))
        if err != nil {
            log.Fatal("HTTP Server error:", err)
        }
    case "client":
        clientConf := viper.Sub("client")
        cert, _ := os.ReadFile(clientConf.GetString("certfile"))
        certPool := x509.NewCertPool()
        certPool.AppendCertsFromPEM(cert)
        tlsConfig := &tls.Config{
            RootCAs: certPool,
        }
        dialer := &websocket.Dialer{
            TLSClientConfig: tlsConfig,
        }
        authConfig := viper.Sub("auth")
        tokenRetriever := auth.TokenRetrieverFunc(auth.RetrieveWithClientSecret)
        token := tokenRetriever.Retrieve(authConfig)
        client := client.NewClient(
            clientConf.GetString("serverAddr"),
            client.Dialer(dialer))
        client.Connect(token)
    case "login":
        authConfig := viper.Sub("auth")
        tokenRetriever := auth.TokenRetrieverFunc(auth.RetrieveWithClientSecret)
        token := tokenRetriever.Retrieve(authConfig)
        log.Print(token)
    }
}
