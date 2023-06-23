package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

type hadlerReciever struct {
	uc usecase.ReceiverMetric
}

func NewReceiver(handler chi.Router, uc usecase.ReceiverMetricRepo) {
	rec := newReceiver(uc)
	//Get /
	handler.Get("/", rec.getMetrics)
	//Get /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	handler.Get(fmt.Sprintf("/value/{%s}/{%s}", _metricType, _metricName), rec.getMetricByName)
	handler.Post("/value/", rec.getMetricByNameJSON)

}

func newReceiver(uc usecase.ReceiverMetricRepo) *hadlerReciever {
	return &hadlerReciever{uc: uc}
}

// getMetricByName принимает запрос ввида /value/{counter|gauge}/someMetric
func (h *hadlerReciever) getMetricByName(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, _metricType)
	metricName := chi.URLParam(r, _metricName)
	result := ""
	var errRepo error
	switch metrics.MetricType(metricType) {
	case metrics.GaugeName:
		gauge, err := h.uc.GetGauge(metricName)
		errRepo = err
		if err != nil {
			break
		}
		result = strconv.FormatFloat(gauge.GetValue(), 'g', -1, 64)
	case metrics.CounterName:
		counter, err := h.uc.GetCounter(metricName)
		errRepo = err
		if err != nil {
			break
		}
		result = strconv.FormatInt(counter.GetValue(), 10)
	}
	if errRepo != nil || result == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = w.Write([]byte(result))
	w.Header().Set("Content-Type", _contentTypeText)
	w.WriteHeader(http.StatusOK)

}

// getMetricByNameJSON принимает запрос ввида /value с body{"ID":"name", "Mtype":"counter|gauge"}
func (h *hadlerReciever) getMetricByNameJSON(w http.ResponseWriter, r *http.Request) {
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

	var errRepo error
	switch metrics.MetricType(m.MType) {
	case metrics.GaugeName:
		gauge, err := h.uc.GetGauge(m.ID)
		errRepo = err
		if err != nil {
			break
		}
		m.Value = new(float64)
		m.SetGauge(gauge)
	case metrics.CounterName:
		counter, err := h.uc.GetCounter(m.ID)
		errRepo = err
		if err != nil {
			break
		}
		m.Delta = new(int64)
		m.SetCounter(counter)
	}
	if errRepo != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	result, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", _contentTypeJSON)
	_, _ = w.Write(result)
	w.WriteHeader(http.StatusOK)
	fmt.Println("cinten result ok", "content", _contentTypeJSON)

}

func (h *hadlerReciever) getMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.uc.GetMetrics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := ""
	for key, value := range metrics {
		result = fmt.Sprintf("%s%s %s\n", result, key, value)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
