package network

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type SocketMessage struct {
	addr    string
	message []byte
}

type Socket struct {
	gossipCh chan SocketMessage

	conn   *websocket.Conn
	server *Server
}

func NewSocket(conn *websocket.Conn, server *Server) *Socket {
	return &Socket{
		gossipCh: make(chan SocketMessage),

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
		senderAddr := []byte(socketMessage.addr + ": ")
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
			addr:    s.conn.RemoteAddr().String(),
			message: payload,
		}
		s.server.broadcastCh <- socketMessage
	}
}
