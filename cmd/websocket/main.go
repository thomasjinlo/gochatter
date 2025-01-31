package main

import (
	"log"

	"github.com/thomasjinlo/gochatter/internal/push"
)

func main() {
	log.Fatal(push.ListenAndServe())
}
