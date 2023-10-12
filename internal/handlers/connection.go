package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/thomasjinlo/gochatter/internal/auth"
	"github.com/thomasjinlo/gochatter/internal/server"
)

func HandleNewConnection(tokenVerifier auth.TokenVerifier, s server.Server) func(http.ResponseWriter, *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }

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

        s.HandleNewConnection(conn, r)
    }
}
