package auth

import "net/http"

type Error struct {
    err error
    Message string
}

func (e Error) Error() string {
    if e.err != nil {
        return e.Message + " " + e.err.Error()
    }

    return e.Message
}

type TokenRetriever interface {
    Retrieve() (string, error)
}

type TokenRetrieverFunc func() (string, error)

func (f TokenRetrieverFunc) Retrieve() (string, error) {
    return f()
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
