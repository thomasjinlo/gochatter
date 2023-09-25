package network

import (
    "fmt"
)

type Server struct {
    broadcastCh chan SocketMessage
    registerCh chan *Socket

    sockets map[*Socket]bool
}

func NewServer() (server *Server) {
    server = &Server{
        broadcastCh: make(chan SocketMessage),
        registerCh: make(chan *Socket),
        sockets: make(map[*Socket]bool),
    }
    go server.handleMessages()

    return server
}

func (s *Server) handleMessages() {
    for {
        select {
        case socket := <-s.registerCh:
            clientIp := socket.conn.RemoteAddr().String()
            fmt.Println(clientIp, "connected")
            s.sockets[socket] = true
        case socketMessage := <-s.broadcastCh:
            for socket := range s.sockets {
                if socket.displayName == socketMessage.displayName {
                    continue
                }

                socket.gossipCh <- socketMessage
            }
        }
    }
}
