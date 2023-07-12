package app

import (
	"context"
	"fmt"
	configSrv "github.com/antoniokichaev/go-alert-me/config/server"
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	postgresRepo "github.com/antoniokichaev/go-alert-me/internal/usecase/repo/postgres"
	"github.com/antoniokichaev/go-alert-me/pkg/memorystorage"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/antoniokichaev/go-alert-me/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func Run() {

	// setup config
	serverConfig := configSrv.NewServerConfig()
	dbConfig := configSrv.NewDBConfig()
	configSrv.ParseFlagServer(serverConfig, dbConfig)

	// setup logger
	l := logger.Initialize("INFO")
	l.Info("config server", zap.Object("config", serverConfig))
	l.Info("config db", zap.Object("config", dbConfig))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//init storage memory as default
	var storage memstorage.Keeper
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
	storage = memstorage.NewMemStorage(storeCounter, storeGauge)

	// we have postgres config we should create it
	db, err := postgres.New(ctx, dbConfig.DatabaseDNS)
	defer func() { _ = db.Close() }()
	if err == nil {
		if err != nil {
			l.Fatal("init db: ", zap.Error(err))
		}

		// create table
		content, err := os.ReadFile("./internal/migrate/postgres/0001_init.sql")
		if err != nil {
			l.Fatal("main ReadFile", zap.Error(err))
		}
		_, err = db.Exec(string(content))

		if err != nil {
			l.Fatal("create tables", zap.Error(err))
		}
		storage = postgresRepo.New(db.DB)
	}

	//create routing
	router := chi.NewRouter()
	router.Use(logger.LogMiddleware)
	router.Use(mgzip.GzipMiddleware)

	{
		updaterUc := usecase.NewUpdater(storage)
		getterUc := usecase.NewReceiver(storage)
		v1.NewRouter(
			router,
			updaterUc,
			getterUc,
			storage,
			l,
		)
	}

	err = http.ListenAndServe(serverConfig.GetMyAddress(), router)
	if err != nil {
		l.Fatal("main: ", zap.String("err", fmt.Sprintf("%v", err)))
	}
}
