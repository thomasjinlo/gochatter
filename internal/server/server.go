package server

import (
    "log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/thomasjinlo/gochatter/internal/client"
)

type Server struct {
    broadcast chan []byte

    clients map[*client.Client]bool
    upgrader websocket.Upgrader
}

func NewHandler() *Server {
    upgrader := websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
    return &Server{
        broadcast: make(chan []byte),
        clients: make(map[*client.Client]bool),
        upgrader: upgrader,
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    conn, _ := s.upgrader.Upgrade(w, r, nil)
    conn.WriteMessage(websocket.BinaryMessage, []byte("hello world"))
    conn.WriteMessage(websocket.BinaryMessage, []byte("YAHOOOOOOOO"))
    log.Print("Connection established!", conn)
    client := client.NewClient(conn, s.broadcast)
    s.clients[client] = true

    go client.ListenWrites()
    go client.ListenReads()
}
