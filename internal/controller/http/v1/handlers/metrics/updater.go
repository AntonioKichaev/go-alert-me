package metrics

import (
	"errors"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/entity"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	_zeroMetricName = http.StatusNotFound
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	_incorrectMetricType         = http.StatusBadRequest
	_incorrectMetricValue        = http.StatusBadRequest
	_contentTypeText             = "text/plain; charset=utf-8"
	_metricType           string = "MetricType"
	_metricName                  = "MetricName"
	_metricValue                 = "MetricValue"
)

type updaterRoutes struct {
	uc usecase.Updater
}

func NewUpdaterRoutes(handler chi.Router, uc usecase.Updater) {
	ur := newUpdaterRoutes(uc)
	handler.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", _metricType, _metricName, _metricValue), ur.updateMetrics)
}
func newUpdaterRoutes(uc usecase.Updater) *updaterRoutes {
	return &updaterRoutes{uc: uc}
}

// updateMetrics принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *updaterRoutes) updateMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, _metricType)
	metricName := chi.URLParam(r, _metricName)
	metricValue := chi.URLParam(r, _metricValue)

	//err := isValidMetricAndMetricName(metricType, metricName)
	var err error

	switch entity.MetricType(metricType) {
	case entity.GaugeName:
		err = h.uc.SetGauge(metricName, metricValue)
	case entity.CounterName:
		err = h.uc.AddCounter(metricName, metricValue)
	default:
		err = ErrorUnknownMetricType
	}
	if err != nil {
		if errors.Is(err, entity.ErrorName) {
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
