package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Auth0Provider struct {
    audience string
    domain string
    issuer string
    secretpath string
}

func NewAuth0Provider(config *viper.Viper) *Auth0Provider {
    return &Auth0Provider{
        audience: config.GetString("aud"),
        domain: config.GetString("domain"),
        issuer: config.GetString("iss"),
        secretpath: config.GetString("secretpath"),
    }
}

func (p Auth0Provider) RetrieveWithClientSecret() (string, error) {
    contentType := "application/x-www-form-urlencoded"

    secret, err := os.ReadFile(p.secretpath)
    if err != nil {
        err := Error{
            Message: "Error reading secret",
            err: err,
        }
        return "", err
    }
    secretStr := strings.TrimRight(string(secret), "\n")

    payload := url.Values{}
    payload.Set("grant_type", "client_credentials")
    payload.Set("audience", p.audience)
    payload.Set("client_id", p.issuer)
    payload.Set("client_secret", secretStr)
    payloadReader := strings.NewReader(payload.Encode())

    res, err := http.Post(p.oauthUrl(), contentType, payloadReader)
    if err != nil {
        err := Error{
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
        err := Error{
            Message: "Error parsing response body",
            err: err,
        }
        return "", err
    }

    return body["access_token"].(string), nil
}

func (p *Auth0Provider) Verify(req *http.Request) error {
    res, err := http.Get(p.jwksUrl())
    if err != nil {
        return Error{
            Message: "Error retrieving JWKs",
            err: err,
        }
    }

    if res.StatusCode != http.StatusOK {
        return Error{
            Message: "Error retrieving JWKs",
        }
    }

    var jwks JsonWebKeySet

    jwksBytes, err := io.ReadAll(res.Body)
    if err != nil {
        return Error{
            Message: "Error reading JWKs response body",
            err: err,
        }
    }

    err = json.Unmarshal(jwksBytes, &jwks)
    if err != nil {
        return Error{
            Message: "Error unmarshalling JWKs",
            err: err,
        }
    }

    parser := jwt.NewParser(
        // issuer of the jwt is Auth0 provider, not the gochatter client
        jwt.WithIssuer(p.url()),
        jwt.WithAudience(p.audience),
        jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
    )

    token, err := parser.Parse(req.Header.Get("Authorization"), func(t *jwt.Token) (interface{}, error) {
        var publicKey *rsa.PublicKey

        for _, jwk := range jwks.Keys {
            if t.Header["kid"] == jwk.KeyId {
                modulus, err := base64.RawURLEncoding.DecodeString(jwk.Modulus)
                if err != nil {
                    return nil, Error{
                        Message: "Error decoding base64 raw url of the modulus",
                        err: err,
                    }
                }

                exponent, err := base64.RawURLEncoding.DecodeString(jwk.Exponent)
                if err != nil {
                    return nil, Error{
                        Message: "Error decoding base64 raw url of the exponent",
                        err: err,
                    }
                }

                publicKey = &rsa.PublicKey{
                    N: big.NewInt(0).SetBytes(modulus),
                    E: int(big.NewInt(0).SetBytes(exponent).Uint64()),
                }
            }
        }

        return publicKey, nil
    })

    if !token.Valid {
        return err
    }

    return nil
}

func (p *Auth0Provider) url() string {
    scheme := "https://"
    path := "/"

    return scheme + p.domain + path
}

func (p *Auth0Provider) oauthUrl() string {
    return p.url() + "oauth/token"
}


func (p *Auth0Provider) jwksUrl() string {
    return p.url() + ".well-known/jwks.json"
}
