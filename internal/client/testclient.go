package client

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
    Author string
    Content string
}

func Connect() {
    tlsConfig := &tls.Config{
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }

    dialer := &websocket.Dialer{TLSClientConfig: tlsConfig}
    log.Println("CONNECTING TO WS")
    header := http.Header{}
    header.Set("DisplayName", "TestClient")
    header.Set("ChannelIds", "1,2,3")
    conn, _, err := dialer.Dial("wss://gochatter.app:443/connect", header)
    log.Println("CONNECTING TO WS")
    if err != nil {
        panic(err)
    }

    log.Println("CONNECTED TO WS SERVER")
    go func() {
        for {
            var message Message
            err := conn.ReadJSON(&message)
            if err != nil {
                panic(err)
            }

            log.Println("RECEIVED MESSAGE FROM WS: ", message)
        }
    }()

    scanner := bufio.NewScanner(os.Stdin)
    for {
        scanner.Scan()

        if err := scanner.Err(); err != nil {
            log.Println("Scanner error: ", err)
            os.Exit(1)
        }

        msg := scanner.Text()

        message := Message{
            Author: "TestClient",
            Content: msg,
        }
        messageBytes, err := json.Marshal(message)
        if err != nil {
            panic(err)
        }

        url := "https://gochatter.app:8443/channels/1/messages"
        res, err := http.Post(url, "application/json", bytes.NewBuffer(messageBytes))
        if err != nil {
            panic(err)
        }
        log.Println("RECEIVED RESPONSE", res)
    }
}
