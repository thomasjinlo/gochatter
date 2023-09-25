package network

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type SocketMessage struct {
	displayName string
	message []byte
}

type Socket struct {
	gossipCh chan SocketMessage

    displayName string
	conn   *websocket.Conn
	server *Server
}

func NewSocket(displayName string, conn *websocket.Conn, server *Server) *Socket {
	return &Socket{
		gossipCh: make(chan SocketMessage),

        displayName: displayName,
		conn:   conn,
		server: server,
	}
}

func (s *Socket) ReceiveFromServer() {
	for {
		socketMessage, open := <-s.gossipCh
		if !open {
			return
		}
		senderAddr := []byte(socketMessage.displayName + ": ")
		message := append(senderAddr, socketMessage.message...)
		s.conn.WriteMessage(websocket.BinaryMessage, message)
	}
}

func (s *Socket) ReceiveFromClient() {
	for {
		_, payload, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Println(s.conn.RemoteAddr().String(), "closed connection")
			s.conn.Close()
			return
		}
		socketMessage := SocketMessage{
			displayName: s.displayName,
			message: payload,
		}
		s.server.broadcastCh <- socketMessage
	}
}
