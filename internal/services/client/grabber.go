package client

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
)

const _additionsMetrics = 2

//go:generate mockery  --name Grabber
type Grabber interface {
	GetSnapshot() map[string]string
}

//go:generate mockery  --name Random
type Random interface {
	Int() int
}
type racoon struct {
	random Random
}

func NewRacoon() Grabber {
	return &racoon{
		random: rand.New(rand.NewSource(322)),
	}
}

func (rc *racoon) GetSnapshot() map[string]string {

	snapshot := make(map[string]string, len(_allowGaugeMetric)+_additionsMetrics)
	rc.setGauge(snapshot)
	rc.setAdditionalMetrics(snapshot)
	return snapshot
}

func (rc *racoon) setGauge(metrics map[string]string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	el := reflect.ValueOf(m)
	for i := 0; i < el.NumField(); i++ {
		field := el.Type().Field(i)
		value := el.Field(i)
		if _, ok := _allowGaugeMetric[field.Name]; ok {
			metrics["gauge/"+field.Name] = fmt.Sprintf("%v", value)
		}
	}

}

func (rc *racoon) setAdditionalMetrics(metrics map[string]string) {
	metrics["gauge/RandomValue"] = strconv.Itoa(rc.random.Int())
	//PollCount - если она counter то зачем мне ее инкриментить если мы просто будетм кидать 1 всегда на сервере будет инкримент
	metrics["counter/PollCount"] = "1"
}
