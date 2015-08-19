package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gorilla/context"

	"foodtastechess/server/auth"
)

func authMiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(res rest.ResponseWriter, req *rest.Request) {
		httpReq := req.Request
		u := context.Get(httpReq, auth.ContextKey)

		req.Env["user"] = u
		handler(res, req)
	}
}

var authMiddleware = rest.MiddlewareSimple(authMiddlewareFunc)
