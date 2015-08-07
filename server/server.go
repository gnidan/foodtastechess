package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/graceful.v1"
	"net"
	"net/http"
	"time"

	"foodtastechess/directory"
	"foodtastechess/logger"
	"foodtastechess/user"
)

var (
	log = logger.Log("server")
)

type Server struct {
	server   *graceful.Server
	Api      *chessApi           `inject:"chessApi"`
	Auth     user.Authentication `inject:"auth"`
	Config   ServerConfig        `inject:"serverConfig"`
	StopChan chan bool           `inject:"stopChan"`
}

type ServerConfig struct {
	BindAddress string
	AppSecret   string
	SessionName string
}

func New() *Server {
	return new(Server)
}

func (s *Server) PreProvide(provide directory.Provider) error {
	provide("serverConfig", ServerConfig{
		BindAddress: "0.0.0.0:8181",
		AppSecret:   "secret12345",
		SessionName: "ftc_session",
	})

	err := provide("chessApi", newChessApi())
	if err != nil {
		log.Error(fmt.Sprintf("Could not provide chess API: %v", err))
	}
	return err
}

func (s *Server) PostPopulate() error {
	//s.Api.init()
	return nil
}

/*
func (s *Server) Start() error {
	secureMux := http.NewServeMux()
	secureMux.HandleFunc("/hello", hello)
	secure := negroni.New()
	secure.Use(s.Auth.LoginRequired())
	secure.UseHandler(secureMux)

	n := negroni.New()
	n.Use(sessions.Sessions(
		s.Config.SessionName,
		cookiestore.New([]byte(s.Config.AppSecret)),
	))
	n.Use(s.Auth.Middleware())

	router := http.NewServeMux()
	router.HandleFunc("/", hello)
	router.Handle("/hello", secure)

	n.UseHandler(router)

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
*/

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/auth/callback", s.Auth.CompleteAuthHandler)
	router.Handle("/api", negroni.New(
		s.Auth.LoginRequired(),
		s.Api.handler(),
	))

	n := negroni.Classic()
	n.UseHandler(router)

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
	s.server.Stop(10 * time.Second)
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

func hello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte("hello, world!<a href=\"/hello\">hello</a>"))
}
