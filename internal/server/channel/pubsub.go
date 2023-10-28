package channel

import "log"

type PubSub struct {
    repository *Repository
}

func NewPubSub(repository *Repository) *PubSub {
    return &PubSub{
        repository: repository,
    }
}

func (ps *PubSub) Subscribe(channelId string, conn *Connection) {
    channel := ps.repository.GetChannel(channelId)
    channel.AddConnection(conn)
    ps.PublishNewConnection(channelId, conn)
}

func (ps *PubSub) PublishNewConnection(channelId string, conn *Connection) {
    channel := ps.repository.GetChannel(channelId)
    channel.BroadcastNewConnection(conn)
}

func (ps *PubSub) PublishNewMessage(channelId, displayName, message string) {
    log.Println("PUBLISHING MESSAGE")
    channel := ps.repository.GetChannel(channelId)
    channel.BroadcastNewMessage(displayName, message)
}
