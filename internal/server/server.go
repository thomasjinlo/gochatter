package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/thomasjinlo/gochatter/internal/server/models"
)

func SetupRoutes() *chi.Mux {
    r := chi.NewRouter()
    log.Println("setting up routes")
    r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("in root")
        w.Write([]byte("HELLO"))
    }))

    r.Route("/channels", func(r chi.Router) {
        r.Route("/{channelId}", func(r chi.Router) {
            //r.Use(channelCtx)
            r.Get("/messages", getChannelMessages)
            r.Post("/messages", createChannelMessage)
        })
    })

    log.Println("Creating channel 1")
    db, err := badger.Open(badger.DefaultOptions("/tmp/gochatter"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.Update(func(txn *badger.Txn) error {
        channel := models.Channel{
            Id: "1",
            Name: "Channel1",
            Users: []string{},
        }
        channelBytes, _ := json.Marshal(channel)
        txn.Set([]byte("channel:1"), channelBytes)
        return nil
    })

    return r
}

func channelCtx(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("received channel message 1")
        channelId := chi.URLParam(r, "channelId")
        channel, err := models.GetChannel(channelId)
        if err != nil {
            http.Error(w, http.StatusText(404), 404)
            return
        }

        ctx := context.WithValue(r.Context(), "channel", channel)
        log.Println("received channel message 2")
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

type ChannelMessage struct {
    DisplayName string
    Message string
}

// /channels/:channel_id/messages
func createChannelMessage(w http.ResponseWriter, r *http.Request) {
    log.Println("[CHANNELS]: RECEIVED MESSAGE")
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }

    var channelMessage ChannelMessage
    err = json.Unmarshal(body, &channelMessage)
    if err != nil {
        log.Fatal(err)
    }

    channelId := chi.URLParam(r, "channelId")
    log.Println("GETTING CHANNEL ", channelId)
    channel, err := models.GetChannel(channelId)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("CONNECTING WITH NATS SERVER")
    channel.CreateMessage(channelMessage.DisplayName, channelMessage.Message)
    nc, err := nats.Connect("nats://nats-server:4222")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("PUBLISHING to NATS SERVER")

    data, err := json.Marshal(channelMessage)
    if err != nil {
        log.Fatal(err)
    }
    nc.Publish("channel.1", data)
    log.Println("PUBLISHED")
}

func getChannelMessages(w http.ResponseWriter, r *http.Request) {
}
