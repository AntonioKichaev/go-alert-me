package main

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/storages/memstorage"
	"net/http"
)

func main() {
	mu := http.NewServeMux()
	storeKeeper := memstorage.NewMemStorage()
	handlerCounter := metrics.NewHandlerMetrics(storeKeeper)
	handlerCounter.Register(mu)
	err := http.ListenAndServe(":8080", mu)
	if err != nil {
		panic(fmt.Errorf("main: %v", err))
	}
}
