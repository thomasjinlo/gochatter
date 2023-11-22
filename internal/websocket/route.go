package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()


	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[Gochatter WS Server] - Healthy"))
		w.WriteHeader(http.StatusOK)
	}))

	channelToUsers := make(map[string][]User)
	r.Get("/connect", addUser(channelToUsers))
	go listenToNewMessages(channelToUsers)

	return r
}

type User struct {
	displayName string
	conn *websocket.Conn
}

func addUser(channelToUsers map[string][]User) func(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		channelIds := r.Header.Get("ChannelIds")
		displayName := r.Header.Get("DisplayName")

		user := User{
			displayName: displayName,
			conn: conn,
		}

		log.Println(fmt.Sprintf("Received connection request from %s, for channels %s", displayName, channelIds))

		for _, channelId := range strings.Split(channelIds, ",") {
			channelToUsers[channelId] = append(channelToUsers[channelId], user)
		}
		w.WriteHeader(http.StatusOK)
	}
}

type Message struct {
	Author string
	Content string
}

func listenToNewMessages(channelToUsers map[string][]User) {
	natsUrl := fmt.Sprintf("nats://%s:4222", os.Getenv("NATS_SERVER"))
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		panic(err)
	}

	nc.Subscribe("channel.*", func(m *nats.Msg) {
		subject := strings.Split(m.Subject, ".")
		channelId := subject[1]
		var message Message
		if err := json.Unmarshal(m.Data, &message); err != nil {
			log.Fatal(err)
		}
		log.Println(fmt.Sprintf("[channel.%v] - {%s}: '%s'", channelId, message.Author, message.Content))

		for _, user := range channelToUsers[channelId] {
			user.conn.WriteJSON(message)
		}
	})
}
