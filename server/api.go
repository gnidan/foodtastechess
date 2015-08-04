package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

type chessApi struct {
	restApi *rest.Api
}

func newChessApi() *chessApi {
	api := new(chessApi)
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
	return api
}

func (api *chessApi) handler() http.Handler {
	return api.restApi.MakeHandler()
}

func hello(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson("hello, world")
}
