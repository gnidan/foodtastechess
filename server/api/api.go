package api

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gorilla/context"
	"net/http"

	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/server/auth"
	"foodtastechess/user"
)

var log = logger.Log("chessApi")

type ChessApi struct {
	ClientQueries queries.ClientQueries `inject:"clientQueries"`

	restApi *rest.Api
}

func New() *ChessApi {
	return new(ChessApi)
}

func (api *ChessApi) PostPopulate() error {
	restApi := rest.NewApi()
	restApi.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/games", api.GetGames),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Could not initialize Chess API: %v", err))
		return err
	}

	restApi.SetApp(router)

	api.restApi = restApi

	return nil
}

func (api *ChessApi) Handler() http.Handler {
	return api.restApi.MakeHandler()
}

func (api *ChessApi) GetGames(res rest.ResponseWriter, req *rest.Request) {
	httpReq := req.Request
	u := context.Get(httpReq, auth.ContextKey).(user.User)

	res.WriteJson(fmt.Sprintf("Hello, %s!", u.NickName))
}
