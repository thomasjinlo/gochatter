package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
    "crypto/tls"
	"syscall"
)

const (
    //SERVER_ADDR string = "71.212.159.238:8443"
    SERVER_ADDR string = "192.168.0.14:8443"
)

type CloseChannel chan struct{}

type Client struct {
    host string
    port string
    protocol string

    conn net.Conn
}

func NewClient(host, port, protocol string) *Client {
    return &Client {
        host: host,
        port: port,
        protocol: protocol,
    }
}

func (c *Client) Connect() {
    config := &tls.Config{
        InsecureSkipVerify: true,
    }
    address := c.host + ":" + c.port
    conn, err := tls.Dial(c.protocol, address, config)
    if err != nil {
        log.Fatal(err)
    }

    closeCh := CloseChannel(make(chan struct{}))
    c.conn = conn

    go handleInterrupt(closeCh)
    go handleServerShutdown(conn, closeCh)
    go handleWrites(conn)

    fmt.Println("Successfully connected to Server!")

    <-closeCh
    conn.Close()
}

func handleWrites(conn net.Conn) {
    reader := bufio.NewReader(os.Stdin)
    for {
        messageBytes, err := reader.ReadBytes('\n')
        if err != nil {
            log.Fatal("Could not read string")
            return
        }

        _, err = conn.Write(messageBytes)
        if err != nil {
            fmt.Println("error writing to connection", err)
            return
        }
    }
}

func handleInterrupt(closeCh CloseChannel) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    <-sigCh
    close(closeCh)
}

func handleServerShutdown(conn net.Conn, closeCh CloseChannel) {
    reader := bufio.NewReader(conn)
    for {
        _, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Server closed the connection:", err)
            // check for channel closure. Channel will already be closed when
            // forcefully interrupted.
            select {
            case _, ok := <-closeCh:
                if ok {
                    close(closeCh)
                }
            default:
                close(closeCh)
            }
            return
        }
    }
}
