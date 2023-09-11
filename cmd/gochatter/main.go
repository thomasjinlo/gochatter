package main

import (
	"fmt"
    "os"

    "github.com/thomasjinlo/gochatter/internal/server"
    "github.com/thomasjinlo/gochatter/internal/client"
    "github.com/thomasjinlo/gochatter/internal/utils"
)

const (
    CONN_HOST string = "192.168.0.14"
    CONN_PORT string = "8443"
    CONN_TYPE string = "tcp4"
)

func main() {
    switch os.Args[1] {
    case "server":
        host := utils.GetPrivateIp()
        port := "8443"
        protocol := server.TCP4

        gochatter := server.NewServer(host, port, protocol)
        gochatter.Listen()
    case "client":
        host := "192.168.0.14"
        port := "8443"
        protocol := server.TCP4

        client := client.NewClient(host, port, protocol)
        client.Connect()
    default:
        fmt.Println("In default")
    }
}
