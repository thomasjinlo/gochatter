package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server interface {
    HandleNewConnection(c *websocket.Conn, r *http.Request)
    Broadcast(socket Socket, message []byte)
}

type SimpleServer struct {
    socketsCh chan Socket

    sockets map[Socket]bool
}

func (s *SimpleServer) HandleNewConnection(c *websocket.Conn, r *http.Request) {
    displayName := r.Header.Get("Display")
    socket := NewSimpleSocket(displayName, c, s)
    s.socketsCh <- socket
}

func NewSimpleServer() *SimpleServer {
    s := &SimpleServer{
        socketsCh: make(chan Socket),
        sockets: make(map[Socket]bool),
    }
    go s.handleMessages()

    return s
}

func (s *SimpleServer) Broadcast(sender Socket, message []byte) {
    for socket := range s.sockets {
        if socket.Identifier() == sender.Identifier() {
            continue
        }
        socket.SendMessage(sender, message)
    }
}

func (s *SimpleServer) handleMessages() {
    for {
        select {
        case socket := <-s.socketsCh:
            fmt.Println(socket.Addr(), "connected")
            s.sockets[socket] = true
        }
    }
}
