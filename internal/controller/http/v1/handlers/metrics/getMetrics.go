package metrics

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/entity"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/go-chi/chi/v5"
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

}

func newReceiver(uc usecase.ReceiverMetricRepo) *hadlerReciever {
	return &hadlerReciever{uc: uc}
}

// HandlerCounter принимает запрос ввида /value/{counter|gauge}/someMetric
func (h *hadlerReciever) getMetricByName(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, _metricType)
	metricName := chi.URLParam(r, _metricName)
	result := ""
	var errRepo error
	switch entity.MetricType(metricType) {
	case entity.GaugeName:
		gauge, err := h.uc.GetGauge(metricName)
		errRepo = err
		if err != nil {
			break
		}
		result = strconv.FormatFloat(gauge.GetValue(), 'g', -1, 64)
	case entity.CounterName:
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
