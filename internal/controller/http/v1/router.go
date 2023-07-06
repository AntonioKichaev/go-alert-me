package v1

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1/handlers/database"
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1/handlers/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux,
	updater usecase.Updater,
	receiver usecase.ReceiverMetric,
	databaseUnit database.StorageStatus,
) {
	up := metrics.NewUpdaterRoutes(updater)
	handler.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", metrics.MetricType, metrics.MetricName, metrics.MetricValue), up.UpdateMetrics)
		r.Post("/", up.UpdateMetricsJSON)
	})

	rec := metrics.NewReceiver(receiver)

	handler.Get("/", rec.GetMetrics)
	handler.Route("/value", func(r chi.Router) {
		r.Post("/", rec.GetMetricByNameJSON)
		//Get /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
		r.Get(fmt.Sprintf("/{%s}/{%s}", metrics.MetricType, metrics.MetricName), rec.GetMetricByName)

	})

	dbHandlers := database.NewHandlers(databaseUnit)
	handler.Get("/ping", dbHandlers.GetPing)

}
