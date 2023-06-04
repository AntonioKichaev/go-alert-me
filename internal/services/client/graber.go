package client

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
)
//go:generate mockery  --name Grabber
type Grabber interface {
	GetSnapshot() map[string]string
}

type Racoon struct {
	PollCount int
}

func NewRacoon() Grabber {
	return &Racoon{}
}

func (rc *Racoon) GetSnapshot() map[string]string {
	snap := rc.getGauge()
	return snap
}

func (rc *Racoon) getGauge() map[string]string {
	const additionsMetrics = 2
	gauge := make(map[string]string, len(_allowGaugeMetric)+additionsMetrics)
	rc.PollCount++
	gauge["PollCount"] = strconv.Itoa(rc.PollCount)
	gauge["RandomValue"] = strconv.Itoa(rand.Int())
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	el := reflect.ValueOf(m)
	for i := 0; i < el.NumField(); i++ {
		field := el.Type().Field(i)
		value := el.Field(i)
		if _, ok := _allowGaugeMetric[field.Name]; ok {
			gauge[field.Name] = fmt.Sprintf("%v", value)
		}
	}
	return gauge
}
