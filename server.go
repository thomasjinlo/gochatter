package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
    CONN_HOST string = "127.0.0.1"
    CONN_PORT string = "5000"
    CONN_TYPE string = "tcp4"
)

func handleConnection(conn net.Conn, wg *sync.WaitGroup, done chan bool) {
    defer wg.Done()
    defer conn.Close()

    fmt.Println("Connected to client!", conn.RemoteAddr().String())
    <-done
}


func main() {
    ADDR := CONN_HOST + ":" + CONN_PORT
    listener, err := net.Listen(CONN_TYPE, ADDR)
    if err != nil {
        log.Fatal(err)
    }

    defer listener.Close()

    fmt.Println("Waiting for incoming connections...")

    done := make(chan bool)
    var wg sync.WaitGroup

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigCh
        fmt.Println("\nShutting down...")
        listener.Close()
    }()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }
        wg.Add(1)
        go handleConnection(conn, &wg, done)
    }

    wg.Wait()
}
