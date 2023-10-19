package prompt

import (
	"bytes"
	"fmt"
	"os"
)

func DisplayName() (string, error) {
    fmt.Print("Enter display name: ")

    buf := make([]byte, 1024)
    n, err := os.Stdin.Read(buf)
    if err != nil {
        return "", err
    }

    displayNameBytes := bytes.TrimRight(buf[:n], "\n")
    return string(displayNameBytes), nil
}
