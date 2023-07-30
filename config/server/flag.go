package server

import (
	"flag"
	"github.com/caarlos0/env/v8"
)

// ParseFlagServer Аргументы агента:
// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
func ParseFlagServer(srv *Server, config *DBConfig) {
	flag.StringVar(&srv.HTTPServerAdr, "a", "localhost:8080", "where server port wil have started")
	flag.IntVar(&srv.StoreIntervalSecond, "i", 300, "timeout to save data, if 0 every operation")
	flag.StringVar(&srv.FileStoragePath, "f", "/tmp/metrics-db.json", "if nothing without save")
	flag.BoolVar(&srv.Restore, "r", true, "restore data from fileStoragePath")
	flag.StringVar(&config.DatabaseDNS, "d", "", "DNS address to connect database")
	flag.StringVar(&srv.SecretKey, "k", "", "Secret key for sign data")
	flag.Parse()
	err := env.Parse(srv)
	if err != nil {
		panic("errParse env: " + err.Error())
	}
	err = env.Parse(config)
	if err != nil {
		panic("errParse env DB: " + err.Error())
	}
}
