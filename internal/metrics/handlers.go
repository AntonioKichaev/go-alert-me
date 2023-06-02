package metrics

import (
	"errors"
	"github.com/antoniokichaev/go-alert-me/internal/handlers"
	"github.com/antoniokichaev/go-alert-me/internal/storages"
	"net/http"
	"strings"
)

const (
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	_zeroMetricName = http.StatusNotFound
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	_incorrectMetricType  = http.StatusBadRequest
	_incorrectMetricValue = http.StatusBadRequest
)

type handlerCouter struct {
	store storages.Keeper
}

func NewHandlerMetrics(store storages.Keeper) handlers.ExecuteHandler {
	return &handlerCouter{store: store}
}

func (h *handlerCouter) Register(router *http.ServeMux) {
	router.HandleFunc("/update/", h.updateMetrics)
}

// HandlerCounter принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *handlerCouter) updateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType, metricName, value := parseDataFromString(r.URL.Path)
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
		gauge, err := newGauge(metricName, value)
		if err != nil {
			w.WriteHeader(_incorrectMetricValue)
			return
		}
		h.store.SetGauge(gauge.name, gauge.value)
	case _counterName:
		counter, err := newCounter(metricName, value)
		if err != nil {
			w.WriteHeader(_incorrectMetricValue)
			return
		}
		h.store.AddCounter(counter.name, counter.value)

	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

// parseDataFromString Парсит строку вида /update/{metricType}/{metricName}/{value}
func parseDataFromString(s string) (string, string, string) {

	key := 2
	// "/update/{metricType}/{metricName}/{value}
	//"[0]/[1]/[2]			/[3]/		 /[4]
	splitS := strings.Split(s, "/")
	result := func() string {
		if key > len(splitS)-1 {
			return ""
		}
		res := splitS[key]
		key++
		return res

	}

	return result(), result(), strings.ReplaceAll(result(), ",", ".")
}
