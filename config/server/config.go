package server

import (
	"fmt"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	HTTPServerAdr       string `env:"ADDRESS"`
	LoggingLevel        string `env:"LOGGING_LEVEL"`
	StoreIntervalSecond int    `env:"STORE_INTERVAL"`
	FileStoragePath     string `env:"FILE_STORAGE_PATH"`
	Restore             bool   `env:"RESTORE"`
	SecretKey           string `env:"KEY"`
}

func (srv *Server) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("HTTPServerAdr", srv.HTTPServerAdr)
	encoder.AddString("LoggingLevel", srv.LoggingLevel)
	encoder.AddInt("StoreIntervalSecond", srv.StoreIntervalSecond)
	encoder.AddString("FileStoragePath", srv.FileStoragePath)
	encoder.AddBool("Restore", srv.Restore)
	encoder.AddString("SecretKey", srv.SecretKey)
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
