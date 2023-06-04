package metrics

import (
	"errors"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers"
	"github.com/antoniokichaev/go-alert-me/internal/storages"
	"net/http"
	"regexp"
)

const (
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	_zeroMetricName = http.StatusNotFound
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	_incorrectMetricType  = http.StatusBadRequest
	_incorrectMetricValue = http.StatusBadRequest
)

var validPath = regexp.MustCompile(`\/update\/(?P<MetricType>\w+)\/(?P<MetricName>\w+)\/(?P<MetricValue>.*)`)

type handlerMetric struct {
	repo storages.MetricRepository
}

func NewHandlerMetrics(store storages.MetricRepository) handlers.ExecuteHandler {
	return newHandlerMetrics(store)
}
func newHandlerMetrics(store storages.MetricRepository) *handlerMetric {
	return &handlerMetric{repo: store}
}

func (h *handlerMetric) Register(router *http.ServeMux) {
	router.HandleFunc("/update/", h.updateMetrics)
}

// HandlerCounter принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *handlerMetric) updateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("income request", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	params := h.GetParams(r.URL.Path)
	metricType := params[_metricType]
	metricName := params[_metricName]
	metricValue := params[_metricValue]
	err = isValidMetricAndMetricName(metricType, metricName)
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
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

// GetParams Парсит строку вида /update/{metricType}/{metricName}/{value}
// на выход получаем map params
func (h *handlerMetric) GetParams(urlPath string) map[string]string {
	names := validPath.SubexpNames()
	a := validPath.FindStringSubmatch(urlPath)
	mp := make(map[string]string, len(a))
	for key := range a {
		mp[names[key]] = a[key]
	}
	return mp
}
