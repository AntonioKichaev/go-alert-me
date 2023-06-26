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
	mu := chi.NewRouter()
	err := logger.Initialize("INFO")
	logger.Log.Info("config server", zap.Object("config", serverConfig))
	if err != nil {
		panic(err)
	}

	mu.Use(logger.LogMiddleware)
	mu.Use(mgzip.GzipMiddleware)

	storeCounter, storeGauge, err := getStoresCounterGauge(serverConfig)
	if err != nil {
		panic(err)
	}
	storeKeeper := memstorage.NewMemStorage(storeCounter, storeGauge)
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

func getStoresCounterGauge(config *configSrv.Server) (*memorystorage.MemoryStorage, *memorystorage.MemoryStorage, error) {
	storeCounter, err := memorystorage.NewMemoryStorage(
		memorystorage.SetStoreIntervalSecond(config.StoreIntervalSecond),
		memorystorage.SetPathToSaveLoad(config.FileStoragePath),
		memorystorage.WithRestore(config.Restore),
	)
	if err != nil {
		return nil, nil, err
	}
	storeGauge, err := memorystorage.NewMemoryStorage(
		memorystorage.SetStoreIntervalSecond(config.StoreIntervalSecond),
		memorystorage.SetPathToSaveLoad(config.FileStoragePath),
		memorystorage.WithRestore(config.Restore),
	)
	return storeCounter, storeGauge, err
}
