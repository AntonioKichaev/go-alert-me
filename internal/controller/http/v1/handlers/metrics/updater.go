package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const (
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	_zeroMetricName = http.StatusNotFound
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	_incorrectMetricType  = http.StatusBadRequest
	_incorrectMetricValue = http.StatusNotFound
	_contentTypeText      = "text/plain; charset=utf-8"
	_contentTypeJSON      = "application/json"
	_contentTypeHTML      = "text/html"
	_metricType           = "MetricType"
	_metricName           = "MetricName"
	_metricValue          = "MetricValue"
)

type updaterRoutes struct {
	uc usecase.Updater
}

func NewUpdaterRoutes(handler chi.Router, uc usecase.Updater) {
	ur := newUpdaterRoutes(uc)
	handler.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", _metricType, _metricName, _metricValue), ur.updateMetrics)
	handler.Post("/update/", ur.updateMetricsJSON)
}
func newUpdaterRoutes(uc usecase.Updater) *updaterRoutes {
	return &updaterRoutes{uc: uc}
}

// updateMetrics принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *updaterRoutes) updateMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, _metricType)
	metricName := chi.URLParam(r, _metricName)
	metricValue := chi.URLParam(r, _metricValue)
	m, err := metrics.NewMetrics(metricType, metricName, metricValue)
	if err == nil {
		switch metrics.MetricType(m.MType) {
		case metrics.GaugeName:
			g, err := h.uc.SetGauge(m.ID, *m.Value)
			if err == nil {
				m.SetGauge(g)
			}
		case metrics.CounterName:
			c, err := h.uc.AddCounter(m.ID, *m.Delta)
			if err == nil {
				m.SetCounter(c)
			}
		default:
			err = metrics.ErrorUnknownMetricType
		}
	}

	if err != nil {
		if errors.Is(err, metrics.ErrorName) {
			w.WriteHeader(_zeroMetricName)
			return
		}
		if errors.Is(err, metrics.ErrorUnknownMetricType) {
			w.WriteHeader(_incorrectMetricType)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

// updateMetricsJSON принимает запрос ввида /update и body metrics struct
func (h *updaterRoutes) updateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	m := &metrics.Metrics{}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = m.IsValid()
	if err != nil {
		if errors.Is(err, metrics.ErrorName) {
			w.WriteHeader(_zeroMetricName)
			return
		}
		if errors.Is(err, metrics.ErrorUnknownMetricType) {
			w.WriteHeader(_incorrectMetricType)
			return
		}
		if errors.Is(err, metrics.ErrorBadValue) {
			w.WriteHeader(_incorrectMetricValue)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch metrics.MetricType(m.MType) {
	case metrics.GaugeName:
		g, err := h.uc.SetGauge(m.ID, *m.Value)
		if err != nil {
			return
		}
		m.SetGauge(g)
	case metrics.CounterName:
		c, err := h.uc.AddCounter(m.ID, *m.Delta)
		if err != nil {
			return
		}
		m.SetCounter(c)
	}

	result, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", _contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)

}
