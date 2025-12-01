package main

import (
	"log"

	"github.com/andreyxaxa/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
