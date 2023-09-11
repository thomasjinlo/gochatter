package server

import (
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "crypto/tls"
    "syscall"
    "path/filepath"
    "strings"
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
}

func NewServer(host, port, protocol string) *Server {
    return &Server{
        host: host,
        port: port,
        protocol: protocol,
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
        message := buf[:n]
        fmt.Println(string(message))
    }
}

func getPrivateAddr() string {
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    // Iterate through the network interfaces
    for _, iface := range interfaces {
        // Filter out loopback and non-up interfaces
        if strings.HasPrefix(iface.Name, "lo") || (iface.Flags&net.FlagUp == 0) {
            continue
        }

        // Get the addresses associated with the interface
        addrs, err := iface.Addrs()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    return ipnet.IP.String()
                }
            }
        }
    }
    return ""
}
