package config

import (
	"flag"
	"github.com/caarlos0/env/v8"
)

// ParseFlag Аргументы агента:
// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
// Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
func ParseFlag(ag *Agent) {
	flag.StringVar(&ag.HTTPServerAdr, "a", "localhost:8080", "where agent wil send request")
	flag.Int64Var(&ag.ReportIntervalSecond, "r", 10, " agent will send report to server in seconds")
	flag.Int64Var(&ag.PollIntervalSecond, "p", 2, "agent will grab data from machine")
	flag.StringVar(&ag.LoggingLevel, "l", "INFO", "agent log level")
	flag.StringVar(&ag.SecretKey, "k", "", "Secret key for sign data")
	flag.Parse()
	err := env.Parse(ag)
	if err != nil {
		panic("errParse env: " + err.Error())
	}

}
