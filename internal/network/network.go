package network

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)


func NewNetworkServer() func(http.ResponseWriter, *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
    hub := NewHub()

    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Fatal("FAILED TO UPGRADE", err)
        }
        node := NewNode(conn, hub)
        hub.registerCh <- node

        go node.ReceiveFromHub()
        go node.ReceiveFromClient()
    }
}
