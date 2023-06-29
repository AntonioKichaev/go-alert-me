package app

import (
	"fmt"
	configSrv "github.com/antoniokichaev/go-alert-me/config/server"
	v1 "github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/antoniokichaev/go-alert-me/pkg/memorystorage"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func Run() {

	serverConfig := configSrv.NewServerConfig()
	configSrv.ParseFlag(serverConfig)
	logger.Initialize("INFO")
	logger.Log.Info("config server", zap.Object("config", serverConfig))

	storeCounter, err := memorystorage.NewMemoryStorage(
		memorystorage.SetStoreIntervalSecond(serverConfig.StoreIntervalSecond),
		memorystorage.SetPathToSaveLoad(serverConfig.FileStoragePath),
		memorystorage.WithRestore(serverConfig.Restore),
	)
	if err != nil {
		logger.Log.Fatal("init storeCounter: ", zap.Error(err))
	}
	storeGauge, err := memorystorage.NewMemoryStorage(
		memorystorage.SetStoreIntervalSecond(serverConfig.StoreIntervalSecond),
		memorystorage.SetPathToSaveLoad(serverConfig.FileStoragePath+".gauge"),
		memorystorage.WithRestore(serverConfig.Restore),
	)
	if err != nil {
		logger.Log.Fatal("init storeGauge: ", zap.Error(err))
	}

	router := chi.NewRouter()
	router.Use(logger.LogMiddleware)
	router.Use(mgzip.GzipMiddleware)

	storeKeeper := memstorage.NewMemStorage(storeCounter, storeGauge)
	{
		updaterUc := usecase.NewUpdater(storeKeeper)
		getterUc := usecase.NewReceiver(storeKeeper)
		v1.NewRouter(router, updaterUc, getterUc)
	}

	err = http.ListenAndServe(serverConfig.GetMyAddress(), router)
	if err != nil {
		logger.Log.Fatal("main: ", zap.String("err", fmt.Sprintf("%v", err)))
	}
}
