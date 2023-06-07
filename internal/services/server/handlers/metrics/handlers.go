package metrics

import (
	"errors"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers"
	"github.com/antoniokichaev/go-alert-me/internal/storages"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	_zeroMetricName = http.StatusNotFound
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	_incorrectMetricType  = http.StatusBadRequest
	_incorrectMetricValue = http.StatusBadRequest
	_contentTypeText      = "text/plain; charset=utf-8"
)

type handlerMetric struct {
	repo storages.MetricRepository
}

func NewHandlerMetrics(store storages.MetricRepository) handlers.ExecuteHandler {
	return newHandlerMetrics(store)
}
func newHandlerMetrics(store storages.MetricRepository) *handlerMetric {
	return &handlerMetric{repo: store}
}

func (h *handlerMetric) Register(router *chi.Mux) {
	router.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", _metricType, _metricName, _metricValue), h.updateMetrics)
	})
}

// updateMetrics принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *handlerMetric) updateMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, _metricType)
	metricName := chi.URLParam(r, _metricName)
	metricValue := chi.URLParam(r, _metricValue)
	err := isValidMetricAndMetricName(metricType, metricName)
	if err != nil {
		if errors.Is(err, ErrorName) {
			w.WriteHeader(_zeroMetricName)
			return
		}
		if errors.Is(err, ErrorUnknownMetricType) {
			w.WriteHeader(_incorrectMetricType)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch MetricType(metricType) {
	case _gaugeName:
		gauge, err := newGauge(metricName, metricValue)
		if err != nil {
			w.WriteHeader(_incorrectMetricValue)
			return
		}
		h.repo.SetGauge(gauge.name, gauge.value)
	case _counterName:
		counter, err := newCounter(metricName, metricValue)
		if err != nil {
			w.WriteHeader(_incorrectMetricValue)
			return
		}
		h.repo.AddCounter(counter.name, counter.value)

	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
