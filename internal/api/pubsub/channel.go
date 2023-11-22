package pubsub

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type ChannelMessage struct {
	Author string
	Content string
}

func PublishNewMessage(id, author, content string) error {
	natsUrl := fmt.Sprintf("nats://%s:4222", os.Getenv("NATS_SERVER"))
	nc, err := nats.Connect(natsUrl)

	sub := fmt.Sprintf("channel.%s", id)
	message := ChannelMessage{
		Author: author,
		Content: content,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return nc.Publish(sub, messageBytes)
}
