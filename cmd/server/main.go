package main

import (
	"fmt"
	configSrv "github.com/antoniokichaev/go-alert-me/config/server"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/storages/memstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	serverConfig := configSrv.NewServerConfig()
	configSrv.ParseFlag(serverConfig)
	mu := chi.NewRouter()
	storeKeeper := memstorage.NewMemStorage()
	handlerKeeper := metrics.NewHandlerMetrics(storeKeeper)
	handlerReciever := metrics.NewHadlerReciever(storeKeeper)

	handlerKeeper.Register(mu)
	handlerReciever.Register(mu)

	err := http.ListenAndServe(serverConfig.GetMyAddress(), mu)
	if err != nil {
		panic(fmt.Errorf("main: %v", err))
	}
}
