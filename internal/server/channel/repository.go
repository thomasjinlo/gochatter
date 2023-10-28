package channel

import (
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
    displayName string
    conn *websocket.Conn
}

func NewConnection(displayName string, conn *websocket.Conn) *Connection {
    return &Connection{
        conn: conn,
        displayName: displayName,
    }
}

type Message struct {
    MessageType string
    Body interface{}
}

func (c *Connection) SendNewConnectionMessage(newConn *Connection) {
    c.conn.WriteJSON(Message{
        MessageType: "NEW_CONNECTION",
        Body: struct {
            NewUser string
        }{
            NewUser: newConn.displayName,
        },
    })
}

func (c *Connection) SendMessage(channelId, sender, message string) {
    log.Println("SENDING JSON MESSAGE")
    c.conn.WriteJSON(Message{
        MessageType: "NEW_MESSAGE",
        Body: struct {
            Sender string
            Message string
            ChannelId string
        }{
            ChannelId: channelId,
            Sender: sender,
            Message: message,
        },
    })
}

type Channel struct {
    id string
    name string
    connections []*Connection
}

func (c *Channel) Id() string {
    return c.id
}

func (c *Channel) Name() string {
    return c.name
}

func (c *Channel) Users() []string {
    var users []string

    for _, conn := range c.connections {
        users = append(users, conn.displayName)
    }

    return users
}

func (c *Channel) BroadcastNewConnection(newConn *Connection) {
    for _, conn := range c.connections {
        if conn == newConn {
            continue
        }

        conn.SendNewConnectionMessage(newConn)
    }
}

func (c *Channel) BroadcastNewMessage(displayName, message string) {
    for _, conn := range c.connections {
        if conn.displayName == displayName {
            continue
        }

        conn.SendMessage(c.id, displayName, message)
    }
}

func (c *Channel) AddConnection(conn *Connection) {
    c.connections = append(c.connections, conn)
}

type Repository struct {
    channels map[string]*Channel
}

func NewRepository() *Repository {
    channels := make(map[string]*Channel)
    channels["1"] = &Channel{id: "1", name: "Channel 1"}
    channels["2"] = &Channel{id: "2", name: "Channel 2"}
    channels["3"] = &Channel{id: "3", name: "Channel 3"}

    return &Repository{
        channels: channels,
    }
}

func (r *Repository) GetChannels() []*Channel {
    var channels []*Channel
    for channelId := range r.channels {
        channels = append(channels, r.GetChannel(channelId))
    }
    return channels
}

func (r *Repository) GetChannel(channelId string) *Channel {
    return r.channels[channelId]
}
