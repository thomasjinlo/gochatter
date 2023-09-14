package network

import (
    "log"
)

type Hub struct {
    broadcastCh chan []byte
    registerCh chan *Node

    nodes map[*Node]bool
}

func NewHub() (hub *Hub) {
    hub = &Hub{
        broadcastCh: make(chan []byte),
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
            log.Print("REGISTERING NODE", node)
            h.nodes[node] = true
        case message := <-h.broadcastCh:
            for node := range h.nodes {
                node.gossipCh <- message
            }
        }
    }
}
