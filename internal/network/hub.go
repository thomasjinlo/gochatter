package network

import (
    "log"
)

type Hub struct {
    broadcastCh chan NodeMessage
    registerCh chan *Node

    nodes map[*Node]bool
}

func NewHub() (hub *Hub) {
    hub = &Hub{
        broadcastCh: make(chan NodeMessage),
        registerCh: make(chan *Node),
        nodes: make(map[*Node]bool),
    }
    go hub.handleMessages()

    return hub
}

func (h *Hub) handleMessages() {
    for {
        select {
        case node := <-h.registerCh:
            clientIp := node.conn.RemoteAddr().String()
            log.Print(clientIp, " connected")
            h.nodes[node] = true
        case nodeMessage := <-h.broadcastCh:
            for node := range h.nodes {
                if node.conn.RemoteAddr().String() == nodeMessage.addr {
                    continue
                }

                node.gossipCh <- nodeMessage
            }
        }
    }
}
