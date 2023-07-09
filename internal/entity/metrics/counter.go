package metrics

import (
	"errors"
	"strconv"
)

type Counter struct {
	Name  string `db:"name"`
	Value int64  `db:"value"`
}

func (c *Counter) GetValue() int64 {
	return c.Value
}
func (c *Counter) SetValue(value int64) {
	c.Value = value
}
func (c *Counter) GetName() string {
	return c.Name
}

func NewCounter(name string, value any) (*Counter, error) {
	name, v, err := validCounter(name, value)
	if err != nil {
		return nil, err
	}
	return &Counter{Name: name, Value: v}, nil
}
func validCounter(name string, value any) (string, int64, error) {
	err := isValidName(name)
	if err != nil {
		return "", 0, ErrorName
	}
	var v int
	switch value := value.(type) {
	case string:
		v, err = strconv.Atoi(value)
	case int:
		v = value
	case int64:
		v = int(value)
	case *int64:
		if value != nil {
			v = int(*value)
		} else {
			err = errors.New("validGauge: nil pointer")
		}

	}

	if err != nil {
		return "", 0, err
	}
	return name, int64(v), err

}
