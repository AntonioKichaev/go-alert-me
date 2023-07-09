package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	metricsEntity "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	_contentTypeText = "text/plain; charset=utf-8"
	_contentTypeJSON = "application/json"
	_contentTypeHTML = "text/html"
	MetricType       = "MetricType"
	MetricName       = "MetricName"
	MetricValue      = "MetricValue"
)

type UpdaterRoutes struct {
	uc     usecase.Updater
	logger *zap.Logger
}

func NewUpdaterRoutes(uc usecase.Updater, logger *zap.Logger) *UpdaterRoutes {
	return &UpdaterRoutes{uc: uc, logger: logger}
}

// UpdateMetrics принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *UpdaterRoutes) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, MetricType)
	metricName := chi.URLParam(r, MetricName)
	metricValue := chi.URLParam(r, MetricValue)
	m, err := metricsEntity.NewMetrics(
		metricsEntity.SetMetricType(metricType),
		metricsEntity.SetName(metricName),
		metricsEntity.SetValueOrDelta(metricValue))

	if err != nil {
		if errors.Is(err, metricsEntity.ErrorName) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, metricsEntity.ErrorUnknownMetricType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.uc.UpdateMetricByParams(r.Context(), m.ID, m.MType, m.GetTmpValue())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

// UpdateMetricJSON принимает запрос ввида /update и body metrics struct
func (h *UpdaterRoutes) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	m := &metricsEntity.Metrics{}
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
		if errors.Is(err, metricsEntity.ErrorName) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, metricsEntity.ErrorUnknownMetricType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, metricsEntity.ErrorBadValue) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err = h.uc.UpdateMetric(r.Context(), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
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

// UpdateMetricsBatchJSON принимает запрос ввида /updates и body metric structs
func (h *UpdaterRoutes) UpdateMetricsBatchJSON(w http.ResponseWriter, r *http.Request) {
	metricsRaw := make([]metricsEntity.Metrics, 0)
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &metricsRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metrics := make([]metricsEntity.Metrics, 0, len(metricsRaw))
	for _, val := range metricsRaw {
		err = val.IsValid()
		if err != nil {
			continue
		}
		metrics = append(metrics, val)
	}

	if err != nil {
		if errors.Is(err, metricsEntity.ErrorName) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, metricsEntity.ErrorUnknownMetricType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, metricsEntity.ErrorBadValue) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(metrics) != 0 {
		err = h.uc.UpdateMetricBatch(r.Context(), metrics)
	}

	if err != nil {
		w.Header().Set("Content-Type", _contentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"err": "%s"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
}
