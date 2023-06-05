package server

import (
	"fmt"
)

type Server struct {
	httpServerAdr string
}

func (srv *Server) GetMyAddress() string {
	return srv.httpServerAdr
}
func (srv *Server) String() string {
	return fmt.Sprintf("server:(%s)", srv.httpServerAdr)
}

func NewServerConfig() *Server {
	return &Server{}
}
