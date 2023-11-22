package models

import (
	"fmt"
	"log"
        "encoding/json"

	badger "github.com/dgraph-io/badger/v4"
)

type Channel struct {
    Id string
    Name string
    Users []string
}

func GetChannel(id string) (*Channel, error) {
    db, err := badger.Open(badger.DefaultOptions("/tmp/gochatter"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var channel Channel

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
        return &channel, err
    }

    return &channel, nil
}

func (c *Channel) CreateMessage(displayName, message string) {
    db, err := badger.Open(badger.DefaultOptions("/tmp/gochatter"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    key := []byte(fmt.Sprintf("channelmessage:%s:", c.Id))
    seq, err := db.GetSequence(key, 1)
    if err != nil {
        log.Fatal(err)
    }
    messageId, err := seq.Next()
    if err != nil {
        log.Fatal(err)
    }
    err = db.Update(func(txn *badger.Txn) error {
        key := []byte(fmt.Sprintf("channelmessage:%s:%v", c.Id, messageId))
        val := []byte(fmt.Sprintf("%s: %s", displayName, message))
        return txn.Set(key, val)
    })
    if err != nil {
        log.Fatal(err)
    }
}
