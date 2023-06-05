package server

import (
	"fmt"
)

type Option func(server *Server)
type Server struct {
	HTTPServerAdr string `env:"ADDRESS"`
}

func (srv *Server) GetMyAddress() string {
	return srv.HTTPServerAdr
}
func (srv *Server) String() string {
	return fmt.Sprintf("server:(%s)", srv.HTTPServerAdr)
}
func SetHttpServerAdr(adr string) Option {
	return func(server *Server) {
		server.HTTPServerAdr = adr
	}
}
func NewServerConfig(opts ...Option) *Server {
	const defaultAdr = "localhost:8080"
	srv := &Server{
		HTTPServerAdr: defaultAdr,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}
