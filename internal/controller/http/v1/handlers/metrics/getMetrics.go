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

type HadlerReciever struct {
	uc usecase.ReceiverMetric
}

func NewReceiver(uc usecase.ReceiverMetric) *HadlerReciever {
	return &HadlerReciever{uc: uc}

}

// GetMetricByName принимает запрос ввида /value/{counter|gauge}/someMetric
func (h *HadlerReciever) GetMetricByName(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, MetricType)
	metricName := chi.URLParam(r, MetricName)
	result := ""
	metric, err := h.uc.GetMetricByName(r.Context(), metricName, metricType)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if metric.Delta != nil {
		result = strconv.FormatInt(*metric.Delta, 10)
	} else {
		result = strconv.FormatFloat(*metric.Value, 'g', -1, 64)
	}

	w.Header().Set("Content-Type", _contentTypeText)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))

}

// GetMetricByNameJSON принимает запрос ввида /value с body{"ID":"name", "Mtype":"counter|gauge"}
func (h *HadlerReciever) GetMetricByNameJSON(w http.ResponseWriter, r *http.Request) {
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

	m, err = h.uc.GetMetricByName(r.Context(), m.ID, m.MType)
	if err != nil {
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

func (h *HadlerReciever) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.uc.GetMetrics(r.Context())
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
