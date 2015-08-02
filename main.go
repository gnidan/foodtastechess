package main

import (
	"os"

	"foodtastechess/server"
)

func main() {
	s := server.New()

	s.Serve("0.0.0.0", os.Getenv("PORT"))
}
