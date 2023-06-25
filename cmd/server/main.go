package main

import (
	"fmt"
	configSrv "github.com/antoniokichaev/go-alert-me/config/server"
	v1 "github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	serverConfig := configSrv.NewServerConfig()
	configSrv.ParseFlag(serverConfig)
	mu := chi.NewRouter()
	err := logger.Initialize("INFO")
	logger.Log.Info("config server", zap.Object("config", serverConfig))
	if err != nil {
		panic(err)
	}

	mu.Use(logger.LogMiddleware)
	mu.Use(mgzip.GzipMiddleware)
	storeKeeper := memstorage.NewMemStorage()
	{
		updaterUc := usecase.NewUpdater(storeKeeper)
		getterUc := usecase.NewReceiver(storeKeeper)
		v1.NewRouter(mu, updaterUc, getterUc)
	}

	err = http.ListenAndServe(serverConfig.GetMyAddress(), mu)
	if err != nil {
		logger.Log.Fatal("main: ", zap.String("err", fmt.Sprintf("%v", err)))
	}
}
