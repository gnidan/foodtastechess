package api

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"math/rand"
	"net/http"
	"strconv"

	"foodtastechess/commands"
	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/users"
)

var log = logger.Log("chessApi")

type ChessApi struct {
	Queries  queries.ClientQueries `inject:"clientQueries"`
	Commands commands.Commands     `inject:"commands"`

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
		rest.Get("/games/:id/", api.GetGameInfo),
		rest.Get("/games/:id/history", api.GetGameHistory),
		rest.Get("/games/:id/validmoves", api.GetGameValidMoves),

		rest.Post("/games/create", api.PostCreateGame),
		rest.Post("/games/:id/join", api.PostJoinGame),
		rest.Post("/games/:id/move", api.PostMove),
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
	return http.StripPrefix("/api", api.restApi.MakeHandler())
}

func (api *ChessApi) GetGames(res rest.ResponseWriter, req *rest.Request) {
	u := getUser(req)
	res.WriteJson(api.Queries.UserGames(u.Uuid))
}

func (api *ChessApi) GetGameInfo(res rest.ResponseWriter, req *rest.Request) {
	id := req.PathParam("id")
	intId, err := strconv.Atoi(id)
	gameId := game.Id(intId)
	if err != nil {
		log.Debug("Recieved an invalid gameid, it was not an int: %s", id)
		rest.NotFound(res, req)
		return
	}

	gameInfo, found := api.Queries.GameInformation(gameId)

	if !found {
		log.Debug("Recieved an invalid gameid, it was not an int: %s", id)
		rest.NotFound(res, req)
		return
	}

	res.WriteJson(gameInfo)
}

func (api *ChessApi) GetGameHistory(res rest.ResponseWriter, req *rest.Request) {
	id := req.PathParam("id")
	intId, err := strconv.Atoi(id)
	gameId := game.Id(intId)
	if err != nil {
		log.Debug("Recieved an invalid gameid, it was not an int: %s", id)
		rest.NotFound(res, req)
		return
	}

	history, found := api.Queries.GameHistory(gameId)
	if !found {
		rest.NotFound(res, req)
		return
	}

	res.WriteJson(history)
}

func (api *ChessApi) GetGameValidMoves(res rest.ResponseWriter, req *rest.Request) {
	id := req.PathParam("id")
	intId, err := strconv.Atoi(id)
	gameId := game.Id(intId)
	if err != nil {
		log.Debug("Recieved an invalid gameid, it was not an int: %s", id)
		rest.NotFound(res, req)
	}

	validMoves, found := api.Queries.ValidMoves(gameId)
	if !found {
		rest.NotFound(res, req)
		return
	}

	res.WriteJson(validMoves)
}

func (api *ChessApi) PostCreateGame(res rest.ResponseWriter, req *rest.Request) {
	user := getUser(req)

	type createBody struct {
		Color game.Color `json:"Color"`
	}

	body := new(createBody)
	err := req.DecodeJsonPayload(body)

	if err != nil || body.Color == "" {
		idx := rand.Perm(2)[0]
		body.Color = []game.Color{game.White, game.Black}[idx]
	}

	ok, msg := api.Commands.ExecCommand(
		commands.CreateGame, user.Uuid, map[string]interface{}{
			"color": body.Color,
		},
	)

	if ok {
		res.WriteHeader(http.StatusAccepted)
		res.WriteJson("ok")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		res.WriteJson(map[string]string{"error": msg})
	}
}

func (api *ChessApi) PostJoinGame(res rest.ResponseWriter, req *rest.Request) {
	user := getUser(req)

	intId, err := strconv.Atoi(req.PathParam("id"))
	gameId := game.Id(intId)
	if err != nil {
		rest.NotFound(res, req)
	}

	ok, msg := api.Commands.ExecCommand(
		commands.JoinGame, user.Uuid, map[string]interface{}{
			"gameId": gameId,
		},
	)

	if ok {
		res.WriteHeader(http.StatusAccepted)
		res.WriteJson("ok")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		res.WriteJson(map[string]string{"error": msg})
	}
}

func (api *ChessApi) PostMove(res rest.ResponseWriter, req *rest.Request) {
	user := getUser(req)

	intId, err := strconv.Atoi(req.PathParam("id"))
	gameId := game.Id(intId)
	if err != nil {
		rest.NotFound(res, req)
	}

	type moveBody struct {
		Move game.AlgebraicMove `json:"Move"`
	}

	body := new(moveBody)
	err = req.DecodeJsonPayload(body)
	if err != nil || body.Move == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.WriteJson(map[string]string{"error": "Move must be a move"})
		return
	}

	ok, msg := api.Commands.ExecCommand(
		commands.Move, user.Uuid, map[string]interface{}{
			"move":   body.Move,
			"gameId": gameId,
		},
	)

	if ok {
		res.WriteHeader(http.StatusAccepted)
		res.WriteJson("ok")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		res.WriteJson(map[string]string{"error": msg})
	}
}

func getUser(req *rest.Request) users.User {
	return req.Env["user"].(users.User)
}
