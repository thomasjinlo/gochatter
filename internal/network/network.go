package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

    "github.com/thomasjinlo/gochatter/internal/auth"
)

type RSAKey struct {
    Kty string `json:"kty"`
    Use string `json:"use"`
    N   string `json:"n"`
    E   string `json:"e"`
    Kid string `json:"kid"`
    Alg string `json:"alg"`
}

func HandleNewConnection(tokenVerifier auth.TokenVerifier) func(http.ResponseWriter, *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
    server := NewServer()

    return func(w http.ResponseWriter, r *http.Request) {
        err := tokenVerifier.Verify(r)
        if err != nil {
            log.Fatal("Failed to verify token")
        }

        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Fatal("FAILED TO UPGRADE", err)
        }
        socket := NewSocket(conn, server)
        server.registerCh <- socket

        go socket.ReceiveFromServer()
        go socket.ReceiveFromClient()
    }
}
