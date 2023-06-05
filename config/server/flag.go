package server

import "flag"

// ParseFlag Аргументы агента:
// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
func ParseFlag(srv *Server) {
	flag.StringVar(&srv.HttpServerAdr, "a", "localhost:8080", "where server port wil have started")
	flag.Parse()
}
