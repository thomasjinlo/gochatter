package client

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    conn *websocket.Conn
    broadcastCh chan []byte
}

func NewClient(conn *websocket.Conn, broadcastCh chan []byte) *Client {
    return &Client{
        conn: conn,
        broadcastCh: broadcastCh,
    }
}

func (c *Client) ListenWrites() {
}

func (c *Client) ListenReads() {

}
