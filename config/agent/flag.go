package config

import (
	"flag"
)

// ParseFlag Аргументы агента:
// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
// Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
func ParseFlag(ag *Agent) {
	flag.StringVar(&ag.httpServerAdr, "a", "localhost:8080", "where agent wil send request")
	flag.Int64Var(&ag.reportIntervalSecond, "r", 10, " agent will send report to server in seconds")
	flag.Int64Var(&ag.pollIntervalSecond, "p", 2, "agent will grab data from machine")
	flag.Parse()
}
