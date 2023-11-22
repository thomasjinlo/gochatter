package channels

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type Message struct {
	Author string
	Content string
}

type Channel struct {
	Id string
	Name string
}

func GetById(id string) (Channel, error) {
	var channel Channel

	db, err := badger.Open(badger.DefaultOptions("/tmp/gochatter"))
	if err != nil {
		return channel, err
	}
	defer db.Close()
	err = db.View(func(txn *badger.Txn) error {
		k := []byte(fmt.Sprintf("channel:%s", id))
		item, err := txn.Get(k)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &channel)
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return channel, err
	}

	return channel, nil
}

func GetMessages(id string) ([]Message, error) {
	var messages []Message
	db, err := badger.Open(badger.DefaultOptions("/tmp/gochatter"))
	if err != nil {
		return messages, err
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(fmt.Sprintf("channel:%s:message", id))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var message Message
			item := it.Item()
			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &message)
			})
			if err != nil {
				return err
			}

			messages = append(messages, message)
		}
	  return nil
	})
	if err != nil {
		return messages, err
	}

	return messages, nil
}

func CreateChannelMessage(id, author, content string) error {
	db, err := badger.Open(badger.DefaultOptions("/tmp/gochatter"))
	if err != nil {
		return err
	}
	defer db.Close()

	k := []byte(fmt.Sprintf("channel:%s", id))
	seq, err := db.GetSequence(k, 1)
	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		messageId, err := seq.Next()
		if err != nil {
			return err
		}
		k = []byte(fmt.Sprintf("channel:%s:message:%v", id, messageId))
		message := Message{
			Author: author,
			Content: content,
		}
		messageBytes, err := json.Marshal(message)
		if err != nil {
			return err
		}
		return txn.Set(k, messageBytes)
	})
	if err != nil {
		return err
	}

	return nil
}
