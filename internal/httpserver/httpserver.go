package httpserver

import (
	"os"
    "log"
	"net/http"
    "path/filepath"
)

type ListenServer interface {
    ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error
}

type HttpServer struct {
    handler http.Handler
}

func NewHttpServer(handler http.Handler) *HttpServer {
    return &HttpServer{
        handler: handler,
    }
}

func (h *HttpServer) Listen() {
    currDir, _ := os.Getwd()
    certFile := filepath.Join(currDir, ".ssh", "cfcert.pem")
    keyFile := filepath.Join(currDir, ".ssh", "cfkey.pem")

    err := http.ListenAndServeTLS(":443", certFile, keyFile, h.handler)
    if err != nil {
        log.Fatal("HTTP Server error:", err)
    }
}
