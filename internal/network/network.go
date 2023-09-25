package network

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/thomasjinlo/gochatter/internal/auth"
)

func HandleNewConnection(tokenVerifier auth.TokenVerifier) func(http.ResponseWriter, *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
    server := NewServer()

    return func(w http.ResponseWriter, r *http.Request) {
        err := tokenVerifier.Verify(r)
        if err != nil {
            fmt.Println("Failed to verify token", err)
            return
        }

        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            fmt.Println("FAILED TO UPGRADE", err)
            return
        }

        displayName := r.Header.Get("Display")
        socket := NewSocket(displayName, conn, server)
        server.registerCh <- socket

        go socket.ReceiveFromServer()
        go socket.ReceiveFromClient()
    }
}
