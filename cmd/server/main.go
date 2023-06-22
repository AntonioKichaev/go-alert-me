package main

import (
	"fmt"
	configSrv "github.com/antoniokichaev/go-alert-me/config/server"
	v1 "github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	serverConfig := configSrv.NewServerConfig()
	configSrv.ParseFlag(serverConfig)
	fmt.Println("config server", serverConfig)
	mu := chi.NewRouter()
	storeKeeper := memstorage.NewMemStorage()
	{
		updaterUc := usecase.NewUpdater(storeKeeper)
		getterUc := usecase.NewReceiver(storeKeeper)
		v1.NewRouter(mu, updaterUc, getterUc)
	}

	err := http.ListenAndServe(serverConfig.GetMyAddress(), mu)
	if err != nil {
		panic(fmt.Errorf("main: %v", err))
	}
}
