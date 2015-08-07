package server

import (
	"fmt"
	"github.com/hydrogen18/stoppableListener"
	"net"
	"net/http"
	"os"

	"foodtastechess/directory"
	"foodtastechess/logger"
)

var (
	log = logger.Log("server")
)

type Server struct {
	listener *stoppableListener.StoppableListener
	Api      *chessApi `inject:"chessApi"`
}

func New() *Server {
	return new(Server)
}

func (s *Server) PreProvide(provide directory.Provider) error {
	err := provide("chessApi", newChessApi())
	if err != nil {
		log.Error(fmt.Sprintf("Could not provide chess API: %v", err))
	}
	return err
}

func (s *Server) PostPopulate() error {
	s.Api.init()
	return nil
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}

	s.serve("0.0.0.0", port)
	return nil
}

func (s *Server) Stop() error {
	s.listener.Stop()
	return nil
}

func (s *Server) serve(bindAddress string, port string) {
	s.listen(bindAddress, port)

	http.Handle("/", s.Api.handler())

	go http.Serve(s.listener, nil)
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
