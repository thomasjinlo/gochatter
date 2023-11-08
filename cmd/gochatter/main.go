package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"github.com/thomasjinlo/gochatter/internal/api"
	"github.com/thomasjinlo/gochatter/internal/auth"
	"github.com/thomasjinlo/gochatter/internal/client"
	ws "github.com/thomasjinlo/gochatter/internal/websocket"

	"github.com/thomasjinlo/gochatter/internal/client/prompt"
	"github.com/thomasjinlo/gochatter/internal/client/tui"
)

func main() {
    currDir, _ := os.Getwd()
    viper.SetConfigName("config")
    viper.AddConfigPath(currDir + "/")
    err := viper.ReadInConfig()
    if err != nil { // Handle errors reading the config file
        log.Fatal("fatal error config file: %w", err)
    }

    serverConf := viper.Sub("server")
    oidcConfig := viper.Sub("oidc")
    auth0Config := oidcConfig.Sub("auth0")
    auth0Provider := auth.NewAuth0Provider(auth0Config)

    switch os.Args[1] {
    case "api":
        //tokenVerifier := auth.TokenVerifierFunc(auth0Provider.Verify)

        mux := api.SetupRoutes()
        err := http.ListenAndServeTLS(
            serverConf.GetString("port"),
            serverConf.GetString("certFile"),
            serverConf.GetString("keyFile"),
            mux)
        if err != nil {
            log.Fatal("HTTP Server error:", err)
        }
    case "websocket":
        log.Println("RUNNING WS SERVER")
        mux := ws.SetupRoutes()
        err := http.ListenAndServeTLS(
            ":443",
            serverConf.GetString("certFile"),
            serverConf.GetString("keyFile"),
            mux)
        if err != nil {
            log.Fatal("HTTP Server error:", err)
        }
    case "client":
        displayName, err := prompt.DisplayName()
        if err != nil {
            log.Fatal("Failed to get displayName", err)
        }

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
        tokenRetriever := auth.TokenRetrieverFunc(auth0Provider.RetrieveWithClientSecret)
        client := client.NewClient(
            displayName,
            serverConf.GetString("domain_name"),
            dialer,
            tokenRetriever,
        )
        //renderer := tui.RendererFunc(tui.TviewRender)
        renderer := tui.RendererFunc(tui.BubbleTeaRender)
        renderer.Render(client)
    case "test-client":
        client.Connect()
    }
}
