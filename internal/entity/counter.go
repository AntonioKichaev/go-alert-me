package entity

import (
	"strconv"
)

type Counter struct {
	Name  string
	Value int64
}

func (c *Counter) GetValue() int64 {
	return c.Value
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

	}

	if err != nil {
		return "", 0, err
	}
	return name, int64(v), err

}