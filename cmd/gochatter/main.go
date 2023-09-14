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
)

func main() {
    currDir, _ := os.Getwd()
    closeCh := make(chan bool)
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
        dialer := websocket.Dialer{
            TLSClientConfig: tlsConfig,
        }
        conn, _, _ := dialer.Dial("wss://gochatter.app:443", nil)
        defer conn.Close()
        go func() {
            for {
                _, payload, _ := conn.ReadMessage()
                log.Print("RECEIVED FROM HUB", string(payload))
            }
        }()

        go func() {
            for {
                buf := make([]byte, 1024)
                n, _ := os.Stdin.Read(buf)
                err := conn.WriteMessage(websocket.BinaryMessage, buf[:n])
                if err != nil {
                    log.Fatal("ERROR WRITING TO NODE", err)
                }
            }
        }()
        <-closeCh
    default:
        log.Print("In default")
    }
}
