package api

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strconv"

	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/users"
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
		rest.Get("/games/:id", api.GetGameInfo),
		rest.Get("/games/:id/history", api.GetGameHistory),
		rest.Get("/games/:id/validmoves", api.GetGameValidMoves),
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
	res.WriteJson(api.ClientQueries.UserGames(u))
}

func (api *ChessApi) GetGameInfo(res rest.ResponseWriter, req *rest.Request) {
	id := req.PathParam("id")
	intId, err := strconv.Atoi(id)
	gameId := game.Id(intId)
	if err != nil
	{
		log.debug("Recieved an invalid gameid, it was not an int: %s", val)
		rest.Error(res,err.Error(),http.StatusNotFound)
	}
	else
	{
		res.WriteJson(api.ClientQueries.GameInformation(intId))
	}
}

func (api *ChessApi) GetGameHistory(res rest.ResponseWriter, req *rest.Request) {
	id := req.PathParam("id")
	intId, err := strconv.Atoi(id)
	gameId := game.Id(intId)
	if err != nil
	{
		log.debug("Recieved an invalid gameid, it was not an int: %s", val)
		rest.Error(res,err.Error(),http.StatusNotFound)
	}
	else
	{
		res.WriteJson(api.ClientQueries.GameHistory(intId))
	}
}

func (api *ChessApi) GetGameValidMoves(res rest.ResponseWriter, req *rest.Request) {
	id := req.PathParam("id")
	intId, err := strconv.Atoi(id)
	gameId := game.Id(intId)
	if err != nil
	{
		log.debug("Recieved an invalid gameid, it was not an int: %s", val)
		rest.Error(res,err.Error(),http.StatusNotFound)
	}
	else
	{
		res.WriteJson(api.ClientQueries.ValidMoves(intId))
	}
}

func getUser(req *rest.Request) users.User {
	return req.Env["user"].(users.User)
}

