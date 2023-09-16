package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Dialer interface {
    Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}

type Client struct {
    closeCh chan struct{}

    conn *websocket.Conn
    addr string
    dialer Dialer
}

func NewClient(addr string, dialer Dialer) *Client {
    return &Client{
        addr: addr,
        dialer: dialer,
    }

}

func (c *Client) Connect() {
    conn, _, _ := c.dialer.Dial(c.addr, nil)
    c.conn = conn
    closeCh := make(chan struct{})

    defer conn.Close()

    go c.receiveFromServer()
    go c.writeToServer()

    <-closeCh
}

func (c *Client) receiveFromServer() {
    for {
        _, payload, err := c.conn.ReadMessage()
        if err != nil {
            fmt.Println("connection closed from server")
            close(c.closeCh)
        }
        fmt.Println(string(payload))
    }
}

func (c *Client) writeToServer() {
    for {
        buf := make([]byte, 1024)
        n, _ := os.Stdin.Read(buf)
        byteMessage := bytes.TrimRight(buf[:n], "\n")
        err := c.conn.WriteMessage(websocket.BinaryMessage, byteMessage)
        if err != nil {
            close(c.closeCh)
        }
    }
}
