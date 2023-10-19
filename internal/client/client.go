package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/thomasjinlo/gochatter/internal/auth"
)

type Dialer interface {
    Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}

type Client struct {
    CloseCh chan struct{}
    ToTuiCh chan string
    FromTuiCh chan string

    conn *websocket.Conn
    DisplayName string
    serverDomainName string
    dialer Dialer
    tokenRetriever auth.TokenRetriever
}

func NewClient(displayName, serverDomainName string, dialer Dialer, tokenRetriever auth.TokenRetriever) *Client {
    return &Client{
        CloseCh: make(chan struct{}),
        ToTuiCh: make(chan string),
        FromTuiCh: make(chan string),

        DisplayName: displayName,
        serverDomainName: serverDomainName,
        dialer: dialer,
        tokenRetriever: tokenRetriever,
    }

}

func (c *Client) Connect() *ChannelSocket {
    conn, _, err := c.dialer.Dial("wss://gochatter.app:8443/connect", nil)
    if err != nil {
        panic(err)
    }
    channelSocket := &ChannelSocket{
        conn: conn,
        userJoinedCallbacks: make(map[string][]func(message string)),
        messageReceivedCallbacks: make(map[string][]func(sender string, message string)),
    }

    go channelSocket.receiveFromServer()

    return channelSocket
}

// func (c *Client) Connect() {
//     token, err := c.tokenRetriever.Retrieve()
//     if err != nil {
//         log.Fatal("Error retrieving token ", err)
//     }
//     header := http.Header{}
//     header.Set("Authorization", token)
//     header.Set("Display", c.DisplayName)
//     scheme := "wss"
//     addr := scheme + "://" + c.serverDomainName
//     conn, _, _ := c.dialer.Dial(addr, header)
//     c.conn = conn
//     defer c.conn.Close()
// 
//     go c.receiveFromServer()
//     go c.writeToServer()
// 
//     <-c.CloseCh
// }

func (c *Client) receiveFromServer() {
    for {
        _, payload, err := c.conn.ReadMessage()
        if err != nil {
            fmt.Println("connection closed from server")
            close(c.CloseCh)
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
            close(c.CloseCh)
        }
    }
}

type Channel struct {
    Id string
    Name string
    Users []string
}

func (c *Client) GetChannels() []*Channel {
    channelApi := ChannelApi{
        client: http.Client{},
    }
    rawBody := channelApi.GetChannels()
    var channels []*Channel
    for _, rawChannel := range rawBody.Channels {
        channel := &Channel{
            Id: rawChannel.Id,
            Name: rawChannel.Name,
            Users: rawChannel.Users,
        }
        channels = append(channels, channel)
    }
    log.Println("CHANNELS", channels)
    return channels
}

type ChannelApi struct {
    client http.Client
}

type ChannelBody struct {
    Channels []struct {
        Id string
        Name string
        Users []string
    }
}

func (ca *ChannelApi) GetChannels() ChannelBody {
    res, err := http.Get("https://gochatter.app:8443/channels")
    if err != nil {
        log.Fatal(err)
    }
    var rawBody struct {
        Channels []struct {
            Id string
            Name string
            Users []string
        }
    }
    body, err := io.ReadAll(res.Body)
    json.Unmarshal(body, &rawBody)
    log.Println(rawBody)

    return rawBody
}

func (ca *ChannelApi) Join(channelId, displayName string) {
    res := struct {
        DisplayName string
        ChannelIds []string
    } {
        DisplayName: displayName,
        ChannelIds: []string{channelId},
    }

    json.Marshal(&res)
}

func (ca *ChannelApi) Connect() {
}

type ChannelSocket struct {
    conn *websocket.Conn
    messageReceivedCallbacks map[string][]func(sender, message string)
    userJoinedCallbacks map[string][]func(message string)
}

func (c *ChannelSocket) JoinChannel(channelId, displayName string) {
    message := struct {
        MessageType string
        Body interface{}
    } {
        MessageType: "JOIN",
        Body: struct {
            ChannelIds []string
            DisplayName string
        } {
            ChannelIds: []string{channelId},
            DisplayName: displayName,
        },
    }
    messageBytes, _ := json.Marshal(message)

    //c.conn.WriteJSON(message)
    c.conn.WriteMessage(websocket.BinaryMessage, messageBytes)
}

func (c *ChannelSocket) HandleMessageReceive(channelId string, messageReceived func(sender, message string)) {
    c.messageReceivedCallbacks[channelId] = append(c.messageReceivedCallbacks[channelId], messageReceived)
}

func (c *ChannelSocket) HandleUserJoin(channelId string, userJoined func(message string)) {
    c.userJoinedCallbacks[channelId] = append(c.userJoinedCallbacks[channelId], userJoined)
}

func (c *ChannelSocket) Broadcast(channelId, displayName, message string) {
    messageBody := struct {
        MessageType string
        Body interface{}
    } {
        MessageType: "BROADCAST",
        Body: struct {
            ChannelId string
            DisplayName string
            Message string
        } {
            ChannelId: channelId,
            DisplayName: displayName,
            Message: message,
        },
    }
    messageBytes, _ := json.Marshal(messageBody)
    err := c.conn.WriteMessage(websocket.BinaryMessage, messageBytes)
    if err != nil {
        log.Fatal(err)
    }
}

func (c *ChannelSocket) receiveFromServer() {
    var message struct {
        MessageType string
        Body json.RawMessage
    }

    for {
        _, res, _ := c.conn.ReadMessage()
        json.Unmarshal(res, &message)

        switch message.MessageType {
        case "NEW_MESSAGE":
            var body struct {
                ChannelId string
                Sender string
                Message string
            }
            json.Unmarshal(message.Body, &body)

            for _, messageReceived := range c.messageReceivedCallbacks[body.ChannelId] {
                messageReceived(body.Sender, body.Message)
            }
        }
    }
}
