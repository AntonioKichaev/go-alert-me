package metrics

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers"
	"github.com/antoniokichaev/go-alert-me/internal/storages"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type hadlerReciever struct {
	repo storages.MetricRepository
}

func NewHandlerReceiver(repo storages.MetricRepository) handlers.ExecuteHandler {
	return newHandlerReceiver(repo)
}
func newHandlerReceiver(repo storages.MetricRepository) *hadlerReciever {
	return &hadlerReciever{repo: repo}
}

func (h *hadlerReciever) Register(router *chi.Mux) {
	//Get /
	router.Get("/", h.getMetrics)
	//Get /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	router.Get(fmt.Sprintf("/value/{%s}/{%s}", _metricType, _metricName), h.getMetricByName)
}

// HandlerCounter принимает запрос ввида /value/{counter|gauge}/someMetric
func (h *hadlerReciever) getMetricByName(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, _metricType)
	metricName := chi.URLParam(r, _metricName)
	result := ""
	var errRepo error
	switch MetricType(metricType) {
	case _gaugeName:
		value, err := h.repo.GetGauge(metricName)
		errRepo = err
		result = strconv.FormatFloat(value, 'g', -1, 64)

	case _counterName:
		value, err := h.repo.GetCounter(metricName)
		errRepo = err
		result = strconv.FormatInt(value, 10)

	}
	if result == "" || errRepo != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", _contentTypeText)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))

}

func (h *hadlerReciever) getMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.repo.GetMetrics()
	result := ""
	for key, value := range metrics {
		result = fmt.Sprintf("%s%s %s\n", result, key, value)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
