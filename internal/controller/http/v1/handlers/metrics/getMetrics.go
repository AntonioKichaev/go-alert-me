package metrics

import (
	"encoding/json"
	"fmt"
	metricsEntity "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
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
	switch metricsEntity.MetricType(metricType) {
	case metricsEntity.GaugeName:
		gauge, err := h.uc.GetGauge(metricName)
		errRepo = err
		if err != nil {
			break
		}
		result = strconv.FormatFloat(gauge.GetValue(), 'g', -1, 64)
	case metricsEntity.CounterName:
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
	w.Header().Set("Content-Type", _contentTypeText)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))

}

// getMetricByNameJSON принимает запрос ввида /value с body{"ID":"name", "Mtype":"counter|gauge"}
func (h *hadlerReciever) getMetricByNameJSON(w http.ResponseWriter, r *http.Request) {
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

	var errRepo error
	switch metricsEntity.MetricType(m.MType) {
	case metricsEntity.GaugeName:
		gauge, err := h.uc.GetGauge(m.ID)
		errRepo = err
		if err != nil {
			break
		}
		m.Value = new(float64)
		m.SetGauge(gauge)
	case metricsEntity.CounterName:
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
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}

func (h *hadlerReciever) getMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.uc.GetMetrics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := "<table>"
	for key, value := range metrics {
		result = fmt.Sprintf("%s <tr><td>%s</td> <td>%s</td></tr>", result, key, value)
	}
	result += "</table>"
	w.Header().Set("Content-Type", _contentTypeHTML)
	w.WriteHeader(http.StatusOK)
	if len(metrics) != 0 {
		result = fmt.Sprintf("<html><body>%s</body></html>", result)
		_, _ = w.Write([]byte(result))
	}
}
