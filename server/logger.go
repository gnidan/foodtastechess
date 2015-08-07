package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/mgutz/ansi"
	"net/http"
	"time"
)

type ServerLogger struct {
}

func NewLogger() *ServerLogger {
	return new(ServerLogger)
}

func (l *ServerLogger) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	log.Info(getStartLogMsg(res, req))

	start := time.Now()
	next(res, req)
	duration := time.Since(start)

	log.Info(getFinishLogMsg(res, req, duration))
}

func getStartLogMsg(res http.ResponseWriter, req *http.Request) string {
	request := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	userAgent := userAgentColor(fmt.Sprintf("- %s", req.UserAgent()))

	return fmt.Sprintf("%s %s", request, userAgent)
}

func getFinishLogMsg(res http.ResponseWriter, req *http.Request, duration time.Duration) string {
	request := fmt.Sprintf("%s %s", req.Method, req.URL.Path)

	rw := res.(negroni.ResponseWriter)
	status := colorStatus(rw.Status())

	microseconds := durationColor(fmt.Sprintf("%dms", duration.Nanoseconds()/1000000))

	return fmt.Sprintf("%s %v in %v", status, request, microseconds)
}

func colorStatus(statusCode int) string {
	statusCodeStr := fmt.Sprintf("%d", statusCode)
	if statusCode >= 400 && statusCode < 500 {
		return statusColor400(statusCodeStr)
	} else if statusCode >= 500 {
		return statusColor500(statusCodeStr)
	} else {
		return statusColorSuccess(statusCodeStr)
	}
}

var (
	userAgentColor     = ansi.ColorFunc("black+b")
	statusColor500     = ansi.ColorFunc("red")
	statusColor400     = ansi.ColorFunc("yellow")
	statusColorSuccess = ansi.ColorFunc("green")
	durationColor      = ansi.ColorFunc("cyan")
)
