package main

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/storages/memstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	mu := chi.NewRouter()
	storeKeeper := memstorage.NewMemStorage()
	handlerKeeper := metrics.NewHandlerMetrics(storeKeeper)
	handlerReciever := metrics.NewHadlerReciever(storeKeeper)

	handlerKeeper.Register(mu)
	handlerReciever.Register(mu)

	err := http.ListenAndServe(":8080", mu)
	if err != nil {
		panic(fmt.Errorf("main: %v", err))
	}
}
