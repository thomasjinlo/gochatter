package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thomasjinlo/gochatter/internal/server/channel"
)

func SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    repo := channel.NewRepository()
    ps := channel.NewPubSub(repo)

    mux.HandleFunc("/channels", NewChannelHandler(repo))
    mux.HandleFunc("/connect", NewConnectionHandler(ps))

    return mux
}

func NewChannelHandler(repo *channel.Repository) func(w http.ResponseWriter, r *http.Request) {
    type ChannelBody struct {
        Id string
        Name string
        Users []string
    }

    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.TLS)
        log.Println("HANDLING REQUEST - /channels")

        switch r.Method {
        case "GET":
            var body struct {
                Channels []ChannelBody
            }
            var channels []ChannelBody
            for _, channel := range repo.GetChannels() {
                channelBody := ChannelBody{
                    Id: channel.Id(),
                    Name: channel.Name(),
                    Users: channel.Users(),
                }
                channels = append(channels, channelBody)
            }
            body.Channels = channels

            resBody, _ := json.Marshal(body)

            w.Header().Add("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(resBody)
        case "POST":
        // TODO
        }
    }
}

func NewConnectionHandler(ps *channel.PubSub) func(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{}
    type WebsocketMessage struct {
        MessageType string
        Body json.RawMessage
    }

    type JoinMessageBody struct {
        ChannelIds []string
        DisplayName string
    }

    type BroadcastMessageBody struct {
        ChannelId string
        DisplayName string
        Message string
    }

    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("RECEIVING WEBSOCKET MESSAGE")
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println(err)
        }

        for {
            var message WebsocketMessage
            _, messageBytes, err := conn.ReadMessage()
            if err != nil {
               log.Println(err)
                return
            }
            json.Unmarshal(messageBytes, &message)

            switch message.MessageType {
            case "JOIN":
                var body JoinMessageBody
                json.Unmarshal(message.Body, &body)

                for _, channelId := range body.ChannelIds {
                    connection := channel.NewConnection(body.DisplayName, conn)
                    ps.Subscribe(channelId, connection)
                }
            case "BROADCAST":
                var body BroadcastMessageBody
                json.Unmarshal(message.Body, &body)
                ps.PublishNewMessage(body.ChannelId, body.DisplayName, body.Message)
            }
        }
    }
}
