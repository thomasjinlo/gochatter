package network

import (
    "log"
    "github.com/gorilla/websocket"
)

type Node struct {
    gossipCh chan []byte

    conn *websocket.Conn
    hub *Hub
}

func NewNode(conn *websocket.Conn, hub *Hub) *Node {
    return &Node{
        gossipCh: make(chan []byte),

        conn: conn,
        hub: hub,
    }
}

func (n *Node) ReceiveFromHub() {
    for {
        log.Print("RECEIVED FROM HUB")
        message, open := <-n.gossipCh
        if !open {
            return
        }

        n.conn.WriteMessage(websocket.BinaryMessage, message)
    }
}

func (n *Node) ReceiveFromClient() {
    for {
        _, payload, err := n.conn.ReadMessage()
        log.Print("RECEIVED MESSAGE FROM CLIENT", string(payload))
        if err != nil {
            log.Fatal("ERROR FROM RECEIVING CLIENT MESSAGE", err)
        }
        n.hub.broadcastCh <- payload
    }
}
