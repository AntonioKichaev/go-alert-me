package app

import (
	"context"
	"fmt"
	configSrv "github.com/antoniokichaev/go-alert-me/config/server"
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/antoniokichaev/go-alert-me/pkg/memorystorage"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/antoniokichaev/go-alert-me/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Run() {

	serverConfig := configSrv.NewServerConfig()
	dbConfig := configSrv.NewDBConfig()
	configSrv.ParseFlagServer(serverConfig, dbConfig)

	l := logger.Initialize("INFO")
	l.Info("config server", zap.Object("config", serverConfig))
	l.Info("config db", zap.Object("config", dbConfig))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	db, err := postgres.New(ctx, dbConfig.DatabaseDNS)
	defer func() { _ = db.Close() }()
	if err != nil {
		l.Error("init db: ", zap.Error(err))
	}

	storeCounter, err := memorystorage.NewMemoryStorage(
		memorystorage.WithLogger(l),
		memorystorage.SetStoreIntervalSecond(serverConfig.StoreIntervalSecond),
		memorystorage.SetPathToSaveLoad(serverConfig.FileStoragePath),
		memorystorage.WithRestore(serverConfig.Restore),
	)
	if err != nil {
		l.Fatal("init storeCounter: ", zap.Error(err))
	}
	storeGauge, err := memorystorage.NewMemoryStorage(
		memorystorage.WithLogger(l),
		memorystorage.SetStoreIntervalSecond(serverConfig.StoreIntervalSecond),
		memorystorage.SetPathToSaveLoad(serverConfig.FileStoragePath+".gauge"),
		memorystorage.WithRestore(serverConfig.Restore),
	)
	if err != nil {
		l.Fatal("init storeGauge: ", zap.Error(err))
	}

	router := chi.NewRouter()
	router.Use(logger.LogMiddleware)
	router.Use(mgzip.GzipMiddleware)

	storeKeeper := memstorage.NewMemStorage(storeCounter, storeGauge)
	{
		updaterUc := usecase.NewUpdater(storeKeeper)
		getterUc := usecase.NewReceiver(storeKeeper)
		v1.NewRouter(
			router,
			updaterUc,
			getterUc,
			db,
		)
	}

	err = http.ListenAndServe(serverConfig.GetMyAddress(), router)
	if err != nil {
		l.Fatal("main: ", zap.String("err", fmt.Sprintf("%v", err)))
	}
}
