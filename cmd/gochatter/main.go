package main

import (
    "os"
    "log"
    "net/http"
    "path/filepath"
    "crypto/tls"
    "crypto/x509"

    "github.com/gorilla/websocket"

    "github.com/thomasjinlo/gochatter/internal/network"
    "github.com/thomasjinlo/gochatter/internal/client"
)

func main() {
    currDir, _ := os.Getwd()
    switch os.Args[1] {
    case "server":
        networkHandler := network.NewNetworkServer()
        certFile := filepath.Join(currDir, ".ssh", "cfcert.pem")
        keyFile := filepath.Join(currDir, ".ssh", "cfkey.pem")

        err := http.ListenAndServeTLS(":443", certFile, keyFile, http.HandlerFunc(networkHandler))
        if err != nil {
            log.Fatal("HTTP Server error:", err)
        }
    case "client":
        certFile := filepath.Join(currDir, ".ssh", "cfclient.pem")
        cert, _ := os.ReadFile(certFile)
        certPool := x509.NewCertPool()
        certPool.AppendCertsFromPEM(cert)
        tlsConfig := &tls.Config{
            RootCAs: certPool,
        }
        addr := "wss://gochatter.app:443"
        dialer := &websocket.Dialer{
            TLSClientConfig: tlsConfig,
        }
        client := client.NewClient(addr, client.Dialer(dialer))
        client.Connect()
    }
}
