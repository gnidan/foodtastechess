package server

import (
	"fmt"
	"github.com/hydrogen18/stoppableListener"
	"net"
	"net/http"

	"foodtastechess/logger"
)

var (
	log = logger.Log("server")
)

type Server struct {
	listener *stoppableListener.StoppableListener
	api      *chessApi
}

func (s *Server) Serve(bindAddress string, port string) {
	s.api = newChessApi()
	s.listen(bindAddress, port)

	http.Serve(s.listener, s.api.handler())
}

func (s *Server) listen(bindAddress string, port string) {
	listen := fmt.Sprintf("%s:%s", bindAddress, port)
	log.Notice("Listening at %s", listen)
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		panic("Could not bind address")
	}

	s.listener, err = stoppableListener.New(listener)
	if err != nil {
		panic("Could not rebind new listener")
	}
}
