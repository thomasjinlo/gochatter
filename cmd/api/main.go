package main

import (
	"log"

	"github.com/thomasjinlo/gochatter/internal/api"
)

func main() {
	log.Fatal(api.ListenAndServe())
}
