package server

import (
	"fmt"
)

type Option func(server *Server)
type Server struct {
	HttpServerAdr string
}

func (srv *Server) GetMyAddress() string {
	return srv.HttpServerAdr
}
func (srv *Server) String() string {
	return fmt.Sprintf("server:(%s)", srv.HttpServerAdr)
}
func SetHttpServerAdr(adr string) Option {
	return func(server *Server) {
		server.HttpServerAdr = adr
	}
}
func NewServerConfig(opts ...Option) *Server {
	const defaultAdr = "localhost:8080"
	srv := &Server{
		HttpServerAdr: defaultAdr,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}
