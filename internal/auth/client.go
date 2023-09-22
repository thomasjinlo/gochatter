package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/golang-jwt/jwt/v5"
)

type TokenRetrieveError struct {
    err error
    Message string
}

func (e TokenRetrieveError) Error() string {
    return e.Message + " " + e.err.Error()
}

type TokenRetriever interface {
    Retrieve() (string, error)
}

type TokenRetrieverFunc func() (string, error)

func (f TokenRetrieverFunc) Retrieve() (string, error) {
    return f()
}

func RetrieveWithClientSecret(config *viper.Viper) func() (string, error) {
    return func() (string, error) { return retrieveWithClientSecret(config) }
}

func retrieveWithClientSecret(config *viper.Viper) (string, error) {
    authUrl := config.GetString("authurl")
    contentType := "application/x-www-form-urlencoded"

    secret, err := os.ReadFile(config.GetString("clientSecretPath"))
    if err != nil {
        err := TokenRetrieveError{
            Message: "Error reading secret",
            err: err,
        }
        return "", err
    }

    payload := url.Values{}
    payload.Set("grant_type", "client_credentials")
    payload.Set("client_id", config.GetString("iss"))
    payload.Set("client_secret", strings.TrimRight(string(secret), "\n"))
    payload.Set("audience", config.GetString("aud"))
    payloadReader := strings.NewReader(payload.Encode())

    res, err := http.Post(authUrl, contentType, payloadReader)
    if err != nil {
        err := TokenRetrieveError{
            Message: "Error retrieving access token",
            err: err,
        }
        return "", err
    }

	defer res.Body.Close()
    var body map[string]interface{}
    bytes, _ := io.ReadAll(res.Body)

    err = json.Unmarshal(bytes, &body) 
    if err != nil {
        err := TokenRetrieveError{
            Message: "Error parsing response body",
            err: err,
        }
        return "", err
    }

    return body["access_token"].(string), nil
}

type TokenVerifier interface {
    Verify(req *http.Request) error
}

type TokenVerifierFunc func(req *http.Request) error

func (f TokenVerifierFunc) Verify(req *http.Request) error {
    return f(req)
}

type JsonWebKey struct {
    Algorithm string `json:"alg"`
    KeyType string `json:"kty"`
    KeyId string `json:"kid"`
    Use string `json:"use"`
    Exponent string `json:"e"`
    Modulus string `json:"n"`
}

type JsonWebKeySet struct {
    Keys []*JsonWebKey `json:"keys"`
}

func VerifyWithJWKS(req *http.Request) error {
    jwksUrl := "https://dev-hjb4npw6t0kq7cve.us.auth0.com/.well-known/jwks.json"
    res, err := http.Get(jwksUrl)
    if err != nil {
        log.Fatal("Error retrieving JWKS ", err)
    }

    if res.StatusCode != http.StatusOK {
        log.Fatal("Response status not OK. Received: ", res.StatusCode)
    }

    var jwks JsonWebKeySet

    jwksBytes, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatal("Error reading JWKS response body ", err)
    }

    err = json.Unmarshal(jwksBytes, &jwks)
    if err != nil {
        log.Fatal("Error unmarshalling JWKS ", err)
    }

    log.Print(jwks)

    parser := jwt.NewParser(
        jwt.WithIssuer("https://dev-hjb4npw6t0kq7cve.us.auth0.com/"),
        jwt.WithAudience("wss://gochatter.app:443"),
        jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
    )
    token, err := parser.Parse(req.Header.Get("Authorization"), func(t *jwt.Token) (interface{}, error) {
        var publicKey *rsa.PublicKey
        log.Print(t.Header)

        for _, jwk := range jwks.Keys {
            log.Print("PRINTING KID ", jwk.KeyId)
            issuer, _ := t.Claims.GetIssuer()
            log.Print("PRINTING ISSUER: ", issuer)
            if t.Header["kid"] == jwk.KeyId {
                log.Print("IN HERE")
                exponent, err := base64.RawURLEncoding.DecodeString(jwk.Exponent)
                if err != nil {
                    log.Fatal("Error decoding base64 raw url ", err)
                }

                modulus, err := base64.RawURLEncoding.DecodeString(jwk.Modulus)
                if err != nil {
                    log.Fatal("Error decoding base64 raw url ", err)
                }

                publicKey = &rsa.PublicKey{
                    E: int(big.NewInt(0).SetBytes(exponent).Uint64()),
                    N: big.NewInt(0).SetBytes(modulus),
                }
            }
        }
        return publicKey, nil
    })
    if !token.Valid {
        log.Print("IN TOKEN NOT VALID")
        log.Print("Token not valid ", err)
        return err
    }

    return nil
}
