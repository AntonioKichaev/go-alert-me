package server

import (
	"flag"
	"github.com/caarlos0/env/v8"
)

// ParseFlag Аргументы агента:
// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
func ParseFlag(srv *Server) {
	flag.StringVar(&srv.HttpServerAdr, "a", "localhost:8080", "where server port wil have started")
	flag.Parse()
	err := env.Parse(srv)
	if err != nil {
		panic("errParse env: " + err.Error())
	}
}
