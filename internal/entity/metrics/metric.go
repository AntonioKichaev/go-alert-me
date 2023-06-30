package metrics

import (
	"fmt"
	"strings"
)

type Metrics struct {
	ID       string   `json:"id"`              // имя метрики
	MType    string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta    *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value    *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	tmpValue any      // значение метрики до момента проставления
}

func (m *Metrics) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("name: %s ", m.ID))
	s.WriteString(fmt.Sprintf("type: %s ", m.MType))
	if m.Delta != nil {
		s.WriteString(fmt.Sprintf("value: %d", *m.Delta))
	}
	if m.Value != nil {
		s.WriteString(fmt.Sprintf("value: %f", *m.Value))
	}
	return s.String()
}

func NewMetrics(opts ...Option) (*Metrics, error) {
	m := &Metrics{Delta: new(int64), Value: new(float64)}
	for _, opt := range opts {
		opt(m)
	}
	switch MetricType(m.MType) {
	case GaugeName:
		g, err := NewGauge(m.ID, m.tmpValue)
		if err != nil {
			return nil, err
		}
		*m.Value = g.GetValue()
	case CounterName:
		c, err := NewCounter(m.ID, m.tmpValue)
		if err != nil {
			return nil, err
		}
		*m.Delta = c.GetValue()
	default:
		return nil, ErrorUnknownMetricType
	}

	return m, nil
}

func (m *Metrics) SetValue(value float64) {
	m.MType = GaugeName.String()
	m.Delta = nil
	*m.Value = value
}
func (m *Metrics) SetDelta(delta int64) {
	m.MType = CounterName.String()
	m.Value = nil
	*m.Delta = delta
}

func (m *Metrics) IsValid() error {
	if len(m.ID) == 0 {
		return fmt.Errorf("%w metrics.IsValid() ID: %v", ErrorName, m.ID)
	}
	switch MetricType(m.MType) {
	case GaugeName, CounterName:
	default:
		return fmt.Errorf("%w metrics.IsValid() type: %v", ErrorUnknownMetricType, m.MType)

	}

	if m.Value == nil && m.Delta == nil {
		return fmt.Errorf("%w metrics.IsValid() value: %v delta: %v", ErrorBadValue, m.Value, m.Delta)
	}
	return nil

}

func (m *Metrics) ToGauge() (*Gauge, error) {
	return NewGauge(m.ID, m.Value)
}

func (m *Metrics) ToCounter() (*Counter, error) {
	return NewCounter(m.ID, m.Delta)
}
func (m *Metrics) GetTmpValue() any {
	return m.tmpValue
}
