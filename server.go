package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
    CONN_HOST string = "127.0.0.1"
    CONN_PORT string = "5000"
    CONN_TYPE string = "tcp4"
)

type CloseChannel chan struct{}

func main() {
    ADDR := CONN_HOST + ":" + CONN_PORT
    listener, err := net.Listen(CONN_TYPE, ADDR)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    fmt.Println("Waiting for incoming connections...")

    closeCh := make(CloseChannel)
    connCh := make(chan net.Conn, 4)
    go handleInterrupt(closeCh, connCh)
    go acceptConnections(listener, connCh)

    <-closeCh
}

func handleInterrupt(closeCh CloseChannel, connCh chan net.Conn) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    <-sigCh
    go func() {
        for conn := range connCh {
            conn.Close()
        }
    }()
    close(connCh)
    close(closeCh)
}

func acceptConnections(listener net.Listener, connCh chan net.Conn) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }
        connCh <- conn
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    fmt.Println("Accepted connection:", conn.RemoteAddr().String())

    for {
        buf := make([]byte, 2048)
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Closing connection to:", conn.RemoteAddr().String())
            conn.Close()
            return
        }
        message := buf[:n]
        fmt.Println(string(message))
    }
}
