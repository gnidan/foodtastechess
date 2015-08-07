package server

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"net/http"

	"foodtastechess/queries"
	"foodtastechess/user"
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
	router, err := rest.MakeRouter()
	if err != nil {
		log.Fatal(err)
	}
	restApi.SetApp(router)

	api.restApi = restApi
}

func (api *chessApi) handler() negroni.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		u := context.Get(req, authContextKey).(user.User)

		msg := fmt.Sprintf("Hello, %s!", u.NickName)

		log.Debug("Inside Api")
		res.Write([]byte(msg))
		next(res, req)
	}
}
