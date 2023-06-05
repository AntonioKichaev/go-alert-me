package server

import (
	"fmt"
)

type Server struct {
	HttpServerAdr string `env:"ADDRESS"`
}

func (srv *Server) GetMyAddress() string {
	return srv.HttpServerAdr
}
func (srv *Server) String() string {
	return fmt.Sprintf("server:(%s)", srv.HttpServerAdr)
}

func NewServerConfig() *Server {
	return &Server{}
}
