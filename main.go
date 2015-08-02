package main

import (
	"fmt"
	"github.com/facebookgo/inject"
	"os"

	"foodtastechess/server"
)

type App struct {
	httpServer server.Server `inject:""`
}

func main() {
	var g inject.Graph
	var a App

	// Here the Populate call is creating instances of NameAPI &
	// PlanetAPI, and setting the HTTPTransport on both to the
	// http.DefaultTransport provided above:
	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	a.httpServer.Serve("0.0.0.0", port)
}
