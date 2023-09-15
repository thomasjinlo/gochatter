package network

import (
    "fmt"
    "github.com/gorilla/websocket"
)

type NodeMessage struct {
    addr string
    message []byte
}

type Node struct {
    gossipCh chan NodeMessage

    conn *websocket.Conn
    hub *Hub
}

func NewNode(conn *websocket.Conn, hub *Hub) *Node {
    return &Node{
        gossipCh: make(chan NodeMessage),

        conn: conn,
        hub: hub,
    }
}

func (n *Node) ReceiveFromHub() {
    for {
        nodeMessage, open := <-n.gossipCh
        if !open {
            return
        }

        message := fmt.Sprintf("%s: %s", nodeMessage.addr, string(nodeMessage.message))
        n.conn.WriteMessage(websocket.BinaryMessage, []byte(message))
    }
}

func (n *Node) ReceiveFromClient() {
    for {
        _, payload, err := n.conn.ReadMessage()
        if err != nil {
            n.conn.Close()
            return
        }

        nodeMessage := NodeMessage{
            addr: n.conn.RemoteAddr().String(),
            message: payload,
        }
        n.hub.broadcastCh <- nodeMessage
    }
}
