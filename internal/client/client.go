package client

import (
    "bytes"
    "fmt"
    "net/http"
    "os"

    "github.com/gorilla/websocket"

    "github.com/thomasjinlo/gochatter/internal/auth"
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
        closeCh: make(chan struct{}),

        addr: addr,
        dialer: dialer,
    }

}

func (c *Client) Connect(token auth.JwtToken) {
    header := http.Header{}
    header.Set("Authorization", "Bearer "+string(token))
    conn, _, _ := c.dialer.Dial(c.addr, header)
    c.conn = conn
    defer c.conn.Close()

    go c.receiveFromServer()
    go c.writeToServer()

    <-c.closeCh
}

func (c *Client) receiveFromServer() {
    for {
        _, payload, err := c.conn.ReadMessage()
        if err != nil {
            fmt.Println("connection closed from server")
            close(c.closeCh)
            return
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
            fmt.Println("GETS IN HERE AFTER CLOSING")
            close(c.closeCh)
        }
    }
}
