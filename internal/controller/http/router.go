package http

import (
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1/handlers/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux, updater usecase.Updater, receiver usecase.ReceiverMetric) {
	h := handler.Route("/", func(r chi.Router) {
	})
	metrics.NewUpdaterRoutes(h, updater)
	metrics.NewReceiver(h, receiver)
}
