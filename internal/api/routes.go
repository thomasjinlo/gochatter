package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thomasjinlo/gochatter/internal/api/channels"
	"github.com/thomasjinlo/gochatter/internal/api/pubsub"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[Gochatter API Server] - Healthy"))
		w.WriteHeader(http.StatusOK)
	}))

	r.Route("/channels", func(r chi.Router) {
		r.Route("/{channelId}", func(r chi.Router) {
			r.Get("/", getChannel)
			r.Get("/messages", getChannelMessages)
			r.Post("/messages", createChannelMessage)
		})
	})

	return r
}

type ChannelMessageBody struct {
	Items []ChannelMessageItem
}

type ChannelMessageItem struct {
	Author string
	Content string
}

type ChannelBody struct {
	Id string
	Name string
}

func getChannel(w http.ResponseWriter, r *http.Request) {
	channelId := chi.URLParam(r, "channelId")
	channel, err := channels.GetById(channelId)
	if err != nil {
		log.Println("[API Server] - error getting channel", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channelBody := ChannelBody{
		Id: channel.Id,
		Name: channel.Name,
	}

	jsonData, err := json.Marshal(channelBody)
	if err != nil {
		log.Println("[API Server] - error marshsalling json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getChannelMessages(w http.ResponseWriter, r *http.Request) {
	channelId := chi.URLParam(r, "channelId")
	messages, err := channels.GetMessages(channelId)
	if err != nil {
		log.Println("[API Server] - error getting channel messages", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channelMessageBody := ChannelMessageBody{
		Items: []ChannelMessageItem{},
	}

	for _, message := range messages {
		channelMessageBody.Items = append(channelMessageBody.Items, ChannelMessageItem{
			Author: message.Author,
			Content: message.Content,
		})
	}

	jsonData, err := json.Marshal(channelMessageBody)
	if err != nil {
		log.Println("[API Server] - error marshsalling json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

type ChannelMessage struct {
	Author string
	Content string
}

func createChannelMessage(w http.ResponseWriter, r *http.Request) {
	channelId := chi.URLParam(r, "channelId")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[API Server] - error reading channel message body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newMessage ChannelMessage
	err = json.Unmarshal(body, &newMessage)
	if err != nil {
		log.Println("[API Server] - error reading unmarshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = channels.CreateChannelMessage(channelId, newMessage.Author, newMessage.Content)
	if err != nil {
		log.Println("[API Server] - error creating new channel message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = pubsub.PublishNewMessage(channelId, newMessage.Author, newMessage.Content)
	if err != nil {
		log.Println("[API Server] - error creating new channel message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
