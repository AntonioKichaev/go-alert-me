package metrics

import (
	"errors"
	"fmt"
	"strconv"
)

type MetricType string

const (
	_counterName MetricType = "counter"
	_gaugeName   MetricType = "gauge"
)

var _allowedMetricsName = [...]string{"gauge", "counter"}
var ErrorName = errors.New("counter: error metric name")
var ErrorUnknownMetricType = errors.New("counter: unknown metric type")

type counter struct {
	name  string
	value int64
}
type gauge struct {
	name  string
	value float64
}

func newGauge(name string, value string) (*gauge, error) {
	name, v, err := validGauge(name, value)
	if err != nil {
		return nil, err
	}
	return &gauge{name: name, value: v}, nil
}

func validGauge(name, value string) (string, float64, error) {
	err := isValidName(name)
	if err != nil {
		return "", 0, ErrorName
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", 0, err
	}
	return name, v, err

}

func validCounter(name, value string) (string, int64, error) {
	err := isValidName(name)
	if err != nil {
		return "", 0, ErrorName
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return "", 0, err
	}
	return name, int64(v), err

}
func isValidName(name string) error {
	if name == "" {
		return ErrorName
	}
	return nil
}

func newCounter(name, value string) (*counter, error) {
	name, v, err := validCounter(name, value)
	if err != nil {
		return nil, err
	}
	return &counter{name: name, value: v}, nil
}

// isValidMetricAndMetricName Проверяет тип метрики и ее имя на валидность
func isValidMetricAndMetricName(metricType, metricName string) error {
	if metricName == "" {
		return fmt.Errorf("%w metricName (%s)", ErrorName, metricName)
	}
	isAllowed := false
	for _, val := range _allowedMetricsName {
		if metricType == val {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("%w unknownMetric (%s)", ErrorUnknownMetricType, metricType)
	}

	return nil
}
