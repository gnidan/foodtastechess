package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"gopkg.in/tylerb/graceful.v1"
	"net"
	"net/http"
	"time"

	"foodtastechess/directory"
	"foodtastechess/logger"
	"foodtastechess/server/api"
	"foodtastechess/server/auth"
	"foodtastechess/server/session"
)

var (
	log = logger.Log("server")
)

type Server struct {
	server   *graceful.Server
	Api      *api.ChessApi       `inject:"chessApi"`
	Auth     auth.Authentication `inject:"auth"`
	Config   ServerConfig        `inject:"serverConfig"`
	StopChan chan bool           `inject:"stopChan"`
}

func New() *Server {
	return new(Server)
}

func (s *Server) PreProvide(provide directory.Provider) error {
	provide("serverConfig", ServerConfig{
		BindAddress: "0.0.0.0:8181",
	})

	err := provide("chessApi", api.New())
	if err != nil {
		log.Error(fmt.Sprintf("Could not provide chess API: %v", err))
		return err
	}

	err = provide("sessionConfig", session.SessionConfig{
		SessionName: "ftc_session",
		Secret:      "secret_123",
	})
	if err != nil {
		log.Error(fmt.Sprintf("Could not provide session config: %v", err))
	}

	err = provide("auth", auth.New())
	if err != nil {
		log.Error(fmt.Sprintf("Could not provide auth service: %v", err))
	}
	return err
}

func (s *Server) Start() error {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(NewLogger())
	n.UseFunc(s.Auth.LoginRequired)
	n.UseHandler(s.Api.Handler())

	s.server = &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    s.Config.BindAddress,
			Handler: n,
		},
	}

	s.listenAndGoServe()

	return nil
}

func (s *Server) Stop() error {
	log.Notice("Stopping server")
	s.server.Stop(1 * time.Second)
	s.StopChan <- true
	return nil
}

func (s *Server) listenAndGoServe() error {
	log.Notice("Listening at %s", s.Config.BindAddress)
	listener, err := net.Listen("tcp", s.Config.BindAddress)
	if err != nil {
		return err
	}

	go s.server.Serve(listener)
	return nil
}
