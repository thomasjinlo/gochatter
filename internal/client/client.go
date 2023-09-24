package client

import (
    "bytes"
    "fmt"
    "net/http"
    "log"
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
    serverDomainName string
    dialer Dialer
    tokenRetriever auth.TokenRetriever
}

func NewClient(serverDomainName string, dialer Dialer, tokenRetriever auth.TokenRetriever) *Client {
    return &Client{
        closeCh: make(chan struct{}),

        serverDomainName: serverDomainName,
        dialer: dialer,
        tokenRetriever: tokenRetriever,
    }

}

func (c *Client) Connect() {
    token, err := c.tokenRetriever.Retrieve()
    if err != nil {
        log.Fatal("Error retrieving token ", err)
    }
    header := http.Header{}
    header.Set("Authorization", token)
    scheme := "wss"
    addr := scheme + "://" + c.serverDomainName
    conn, _, _ := c.dialer.Dial(addr, header)
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
