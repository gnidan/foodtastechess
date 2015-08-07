package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"

	"foodtastechess/queries"
)

type chessApi struct {
	restApi *rest.Api

	ClientQueries queries.ClientQueries `inject:"clientQueries"`
}

func newChessApi() *chessApi {
	return new(chessApi)
}

func (api *chessApi) init() {
	restApi := rest.NewApi()
	restApi.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", hello),
	)
	if err != nil {
		log.Fatal(err)
	}
	restApi.SetApp(router)

	api.restApi = restApi
}

func (api *chessApi) handler() http.Handler {
	return api.restApi.MakeHandler()
}

func hello(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson("hello, world")
}
