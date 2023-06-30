package server

import (
	"flag"
	"github.com/caarlos0/env/v8"
)

// ParseFlag Аргументы агента:
// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
func ParseFlag(srv *Server) {
	flag.StringVar(&srv.HTTPServerAdr, "a", "localhost:8080", "where server port wil have started")
	flag.IntVar(&srv.StoreIntervalSecond, "i", 300, "timeout to save data, if 0 every operation")
	flag.StringVar(&srv.FileStoragePath, "f", "/tmp/metrics-db.json", "if nothing without save")
	flag.BoolVar(&srv.Restore, "r", true, "restore data from fileStoragePath")
	flag.Parse()
	err := env.Parse(srv)
	if err != nil {
		panic("errParse env: " + err.Error())
	}
}
