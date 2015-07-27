package server

import (
	"fmt"
	"github.com/hydrogen18/stoppableListener"
	"net"
	"net/http"
)

type Server struct {
	listener *stoppableListener.StoppableListener
}

func New() *Server {
	s := new(Server)

	http.HandleFunc("/", hello)

	return s
}

func (s *Server) Serve(bindAddress string, port string) {
	address := fmt.Sprintf("%s:%s", bindAddress, port)
	fmt.Println("listening at ", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("Could not bind address")
	}

	s.listener, err = stoppableListener.New(listener)
	if err != nil {
		panic("Could not rebind new listener")
	}

	http.Serve(s.listener, nil)
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
