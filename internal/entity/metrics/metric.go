package metrics

import (
	"fmt"
	"strings"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
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

func NewMetrics(mType, mName, mValue string) (*Metrics, error) {
	m := &Metrics{ID: mName, MType: mType, Delta: new(int64), Value: new(float64)}
	switch MetricType(mType) {
	case GaugeName:
		g, err := NewGauge(mName, mValue)
		if err != nil {
			return nil, err
		}
		*m.Value = g.GetValue()
	case CounterName:
		c, err := NewCounter(mName, mValue)
		if err != nil {
			return nil, err
		}
		*m.Delta = c.GetValue()
	default:
		return nil, ErrorUnknownMetricType
	}

	return m, nil
}

func (m *Metrics) SetGauge(gauge *Gauge) {
	m.ID = gauge.GetName()
	m.MType = GaugeName.String()
	m.Delta = nil
	*m.Value = gauge.GetValue()
}
func (m *Metrics) SetCounter(counter *Counter) {
	m.ID = counter.GetName()
	m.MType = CounterName.String()
	m.Value = nil
	*m.Delta = counter.GetValue()
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
