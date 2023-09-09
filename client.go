package main

import (
	"log"
	"net"
    "fmt"
)


func main() {
    conn, err := net.Dial("tcp4", "127.0.0.1:5000")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    fmt.Println("Connection to server!")
}
