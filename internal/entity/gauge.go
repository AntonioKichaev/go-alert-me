package entity

import "strconv"

type Gauge struct {
	Name  string
	Value float64
}

func (g *Gauge) GetValue() float64 {
	return g.Value
}
func (g *Gauge) GetName() string {
	return g.Name
}
func NewGauge(name string, value any) (*Gauge, error) {
	name, v, err := validGauge(name, value)
	if err != nil {
		return nil, err
	}
	return &Gauge{Name: name, Value: v}, nil
}

func validGauge(name string, value any) (string, float64, error) {
	err := isValidName(name)
	if err != nil {
		return "", 0, ErrorName
	}
	var v float64
	switch value := value.(type) {
	case string:
		v, err = strconv.ParseFloat(value, 64)
	case float64:
		v = value
	}

	if err != nil {
		return "", 0, err
	}
	return name, v, err
}
