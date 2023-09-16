package auth

import (
	"log"
	"os"
	"time"
    "crypto/x509"
	"encoding/pem"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type JwtToken string

type TokenRetriever interface {
    Retrieve(config *viper.Viper) JwtToken
}

type TokenRetrieverFunc func(config *viper.Viper) JwtToken

func (tr TokenRetrieverFunc) Retrieve(config *viper.Viper) JwtToken {
    return tr(config)
}

func RetrieveWithClientAssertion(config *viper.Viper) JwtToken {
    newUUID := uuid.New()
    keyId, err := os.ReadFile(config.GetString("keyid"))
    if err != nil {
        log.Fatal("Error reading key id")
    }

    claims := jwt.MapClaims{
        "alg": config.GetString("alg"),
        "aud": config.GetString("aud"),
        "kid": string(keyId),
        "iss": config.GetString("iss"),
        "sub": config.GetString("iss"),
        "jti": newUUID.String(),
        "exp": time.Now().Unix() + 36000,
    }

    // Create a new JWT token with the claims
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
