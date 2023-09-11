package server

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "path/filepath"
    "strings"
    "syscall"
)

const (
    TCP  = "tcp"
    TCP4 = "tcp4"
    TCP6 = "tcp6"
)

type CloseChannel chan struct{}

type Server struct {
    host string
    port string
    protocol string
    connections []net.Conn
}

func NewServer(host, port, protocol string) *Server {
    return &Server{
        host: host,
        port: port,
        protocol: protocol,
        connections: []net.Conn{},
    }
}

func (s *Server) Listen() {
    currentDir, _ := os.Getwd()
    certFile := filepath.Join(currentDir, ".ssh", "gochatter.crt")
    keyFile := filepath.Join(currentDir, ".ssh", "gochatter.key")
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Error loading certificates: %v", err)
    }

    config := tls.Config{Certificates: []tls.Certificate{cert}}
    address := s.host + ":" + s.port
    fmt.Println("Listening on port:", address)
    listener, err := tls.Listen(s.protocol, address, &config)
    if err != nil {
        log.Fatalf("Error creating listener: %v", err)
    }

    defer listener.Close()
    fmt.Println("Waiting for incoming connections...")

    closeCh := make(CloseChannel)
    connCh := make(chan net.Conn, 4)
    go handleInterrupt(closeCh, connCh)
    go acceptConnections(listener, connCh, s)

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

func acceptConnections(listener net.Listener, connCh chan net.Conn, s *Server) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }
        connCh <- conn
        s.connections = append(s.connections, conn)
        go handleConnection(conn, s)
    }
}

func handleConnection(conn net.Conn, s *Server) {
    fmt.Println(conn.RemoteAddr().String(), "connection created")

    for {
        buf := make([]byte, 2048)
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println(conn.RemoteAddr().String(), "closing connection")
            conn.Close()
            fmt.Println(conn.RemoteAddr().String(), "connection closed")
            return
        }

        messageBytes := bytes.TrimRight(buf[:n], "\n")

        if messageBytes[0] == '/' {
            commands := strings.Split(string(messageBytes[1:]), " ")
            fmt.Println(commands)
            command := commands[0]
            switch command {
            case "ls":
                fmt.Println("Listing connections...")
                for _, connection := range s.connections {
                    targetConnIp := connection.RemoteAddr().String()
                    sourceConnIp := conn.RemoteAddr().String()

                    if targetConnIp == sourceConnIp {
                        continue
                    }

                    _, err := conn.Write([]byte(targetConnIp))
                    if err != nil {
                        fmt.Println("Error writing", err)
                    }
                    fmt.Println(targetConnIp)
                }
            case "connect":
                //otherConnIp := commands[1:]
                //conn.Write(commands)
            }
        } else {
            fmt.Println(conn.RemoteAddr().String(), string(messageBytes))
        }
    }
}
