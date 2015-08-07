package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
	"time"
)

type ServerLogger struct {
}

func NewLogger() *ServerLogger {
	return &ServerLogger{}
}

func (l *ServerLogger) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	log.Info(fmt.Sprintf("Started %s %s", req.Method, req.URL.Path))

	next(res, req)

	rw := res.(negroni.ResponseWriter)
	log.Info(fmt.Sprintf("Completed %v %s in %v", rw.Status(), http.StatusText(rw.Status()), time.Since(start)))
}
