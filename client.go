package main

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
    SERVER_ADDR string = "71.212.159.238:8443"
)

func main() {
    config := &tls.Config{
        InsecureSkipVerify: true, // Insecure for self-signed certificates
    }
    conn, err := tls.Dial("tcp", SERVER_ADDR, config)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    fmt.Println("Connection to server!")

    closeCh := CloseChannel(make(chan struct{}))
    go handleInterrupt(closeCh)
    go handleServerShutdown(conn, closeCh)
    go handleWrites(conn)

    <-closeCh
}

func handleWrites(conn net.Conn) {
    reader := bufio.NewReader(os.Stdin)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal("Could not read string")
            return
        }

        _, err = conn.Write([]byte(message))
        if err != nil {
            fmt.Println("error writing to connection", err)
            return
        }
    }
}

type CloseChannel chan struct{}

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
