package server

import (
	"fmt"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	HTTPServerAdr string `env:"ADDRESS"`
	LoggingLevel  string `env:"LOGGING_LEVEL"`
}

func (srv *Server) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("HTTPServerAdr", srv.HTTPServerAdr)
	encoder.AddString("LoggingLevel", srv.LoggingLevel)
	return nil
}

func (srv *Server) GetMyAddress() string {
	return srv.HTTPServerAdr
}
func (srv *Server) GetLoggingLevel() string {
	return srv.LoggingLevel
}
func (srv *Server) String() string {
	return fmt.Sprintf("server:(%s)", srv.HTTPServerAdr)
}

func NewServerConfig(opts ...Option) *Server {
	const defaultAdr = "localhost:8080"
	const defaultLoggingLevel = "DEBUG"
	srv := &Server{
		HTTPServerAdr: defaultAdr,
		LoggingLevel:  defaultLoggingLevel,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}
