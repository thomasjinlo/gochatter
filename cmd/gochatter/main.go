package main

import (
    "os"
    "log"
    "net/http"
    "crypto/tls"
    "crypto/x509"

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
        serverConf := viper.Sub("server")
        err := http.ListenAndServeTLS(
            serverConf.GetString("port"),
            serverConf.GetString("certFile"),
            serverConf.GetString("keyFile"),
            http.HandlerFunc(network.HandleNewConnection()))
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
        client := client.NewClient(
            clientConf.GetString("serverAddr"),
            client.Dialer(dialer))
        client.Connect()
    case "login":
        authConfig := viper.Sub("auth")
        tokenRetriever := auth.TokenRetrieverFunc(auth.RetrieveWithClientAssertion)
        token := tokenRetriever.Retrieve(authConfig)
        log.Print(token)
    }
}
