package auth

import (
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"net/http"
    "net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type JwtToken string

type TokenRetriever interface {
    Retrieve(config *viper.Viper) JwtToken
}

type TokenRetrieverFunc func(config *viper.Viper) JwtToken

func (f TokenRetrieverFunc) Retrieve(config *viper.Viper) JwtToken {
    return f(config)
}

func RetrieveWithClientSecret(config *viper.Viper) JwtToken {
    authUrl := config.GetString("authurl")
    contentType := "application/x-www-form-urlencoded"

    secret, err := os.ReadFile(config.GetString("secret"))
    if err != nil {
        log.Fatal("Error reading secret")
    }

    payload := url.Values{}
    payload.Set("grant_type", "client_credentials")
    payload.Set("client_id", config.GetString("iss"))
    payload.Set("client_secret", strings.TrimRight(string(secret), "\n"))
    payload.Set("audience", config.GetString("aud"))
    payloadReader := strings.NewReader(payload.Encode())

    res, err := http.Post(authUrl, contentType, payloadReader)
    if err != nil {
        log.Fatal("Error retrieving access token", err)
    }

	defer res.Body.Close()
    token, _ := io.ReadAll(res.Body)

    return JwtToken(token)
}

func RetrieveWithClientAssertion(config *viper.Viper) JwtToken {
    newUUID := uuid.New()
    keyId, err := os.ReadFile(config.GetString("keyid"))
    if err != nil {
        log.Fatal("Error reading key id")
    }
    kid := strings.TrimRight(string(keyId), "\n")

    claims := jwt.MapClaims{
        "alg": config.GetString("alg"),
        "kid": kid,
        "aud": config.GetString("aud"),
        "iss": config.GetString("iss"),
        "sub": config.GetString("iss"),
        "jti": newUUID.String(),
        "exp": time.Now().Unix() + 3600,
    }

    //fmt.Println("signing with claim", claims)

    //// Create a new JWT token with the claims
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

    // Sign the token with a secret key (replace with your secret)
    keyBytes, err := os.ReadFile(config.GetString("privatekey"))
    if err != nil {
        log.Fatal("Error reading private key")
    }
    block, _ := pem.Decode(keyBytes)
    privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        log.Fatal("Error decoding private key", err)
    }

    tokenString, err := token.SignedString(privateKey)
    if err != nil {
        log.Fatal("Error signing token")
    }

    return JwtToken(tokenString)
}

type TokenVerifier interface {
    Verify(token string) error
}

type TokenVerifierFunc func(token string) error

func (f TokenVerifierFunc) Verify(token string) error {
    return f(token)
}

func VerifyClientAssertion(config *viper.Viper) func(string) error {
    //keyBytes, err := os.ReadFile(config.GetString("publickey"))
    //if err != nil {
    //    log.Fatal("Error reading private key")
    //}
    //block, _ := pem.Decode(keyBytes)
    //publicKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    //if err != nil {
    //    log.Fatal("Error decoding private key", err)
    //}

    //return func(token string) error {
    //    tokenParts := strings.Split(token, ".")
    //    signingString := strings.Join(tokenParts[:2], ".")
    //    sig := []byte(tokenParts[2])
    //    return jwt.SigningMethodRS256.Verify(signingString, sig, publicKey)
    //}
    return nil
}
