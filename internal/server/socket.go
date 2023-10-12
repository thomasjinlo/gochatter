package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Socket interface {
    Addr() string
    Identifier() string
    SendMessage(sender Socket, message []byte)
}

type SimpleSocket struct {
    displayName string
	conn   *websocket.Conn
	server Server
}

func NewSimpleSocket(displayName string, conn *websocket.Conn, server Server) *SimpleSocket {
    s := &SimpleSocket{
        displayName: displayName,
		conn:   conn,
		server: server,
	}

    go s.receiveFromClient()

    return s
}

func (s *SimpleSocket) Addr() string {
    return s.conn.RemoteAddr().String()
}

func (s *SimpleSocket) Identifier() string {
    return s.displayName
}

func (s *SimpleSocket) SendMessage(sender Socket, message []byte) {
    senderAddr := []byte(sender.Identifier() + ": ")
    message = append(senderAddr, message...)
    s.conn.WriteMessage(websocket.BinaryMessage, message)
}


func (s *SimpleSocket) receiveFromClient() {
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Println(s.conn.RemoteAddr().String(), "closed connection")
			s.conn.Close()
			return
		}

		s.server.Broadcast(s, message)
	}
}
