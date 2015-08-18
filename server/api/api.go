package api

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"

	"foodtastechess/logger"
	"foodtastechess/queries"
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
	restApi.Use(
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.JsonIndentMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
	)
	restApi.Use(authMiddleware)
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
	u := getUser(req)
	res.WriteJson(fmt.Sprintf("Hello, %s!", u.Name))
}

func getUser(req *rest.Request) user.User {
	return req.Env["user"].(user.User)
}
