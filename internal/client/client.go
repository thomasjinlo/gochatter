package client

import (
    //"bytes"
    "fmt"
    "net/http"
    "log"
    //"os"

    "github.com/gorilla/websocket"

    "github.com/thomasjinlo/gochatter/internal/auth"
)

type Dialer interface {
    Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}

type Client struct {
    closeCh chan struct{}
    ToTuiCh chan string
    FromTuiCh chan string

    conn *websocket.Conn
    displayName string
    serverDomainName string
    dialer Dialer
    tokenRetriever auth.TokenRetriever
}

func NewClient(displayName, serverDomainName string, dialer Dialer, tokenRetriever auth.TokenRetriever) *Client {
    return &Client{
        closeCh: make(chan struct{}),
        ToTuiCh: make(chan string),
        FromTuiCh: make(chan string),

        displayName: displayName,
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
    header.Set("Display", c.displayName)
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
        c.ToTuiCh <- string(payload)
        //fmt.Println(string(payload))
    }
}

func (c *Client) writeToServer() {
    for {
        //buf := make([]byte, 1024)
        //n, _ := os.Stdin.Read(buf)
        //byteMessage := bytes.TrimRight(buf[:n], "\n")
        message := <-c.FromTuiCh
        err := c.conn.WriteMessage(websocket.BinaryMessage, []byte(message))
        if err != nil {
            close(c.closeCh)
        }
    }
}
