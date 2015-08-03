package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"

	"foodtastechess/logger"
)

var log = logger.Log("api")

func apiHandler() http.Handler {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", hello),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api.MakeHandler()
}

func hello(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(map[string]string{"Body": "Hello World!"})
}
