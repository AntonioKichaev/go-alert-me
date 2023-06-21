package grabbers

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
	random       Random
	allowMetrics map[string]struct{}
}

type Option func(rac *racoon)

func SetAllowMetrics(metrics map[string]struct{}) Option {
	return func(rac *racoon) {
		rac.allowMetrics = metrics
	}
}

func NewRacoon(opts ...Option) Grabber {
	rn := &racoon{
		random: rand.New(rand.NewSource(322)),
	}
	for _, opt := range opts {
		opt(rn)
	}
	return rn
}

func (rc *racoon) GetSnapshot() map[string]string {

	snapshot := make(map[string]string, len(rc.allowMetrics)+_additionsMetrics)
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
		if _, ok := rc.allowMetrics[field.Name]; ok {
			metrics["gauge/"+field.Name] = fmt.Sprintf("%v", value)
		}
	}

}

func (rc *racoon) setAdditionalMetrics(metrics map[string]string) {
	metrics["gauge/RandomValue"] = strconv.Itoa(rc.random.Int())
	//PollCount - если она counter то зачем мне ее инкриментить если мы просто будетм кидать 1 всегда на сервере будет инкримент
	metrics["counter/PollCount"] = "1"
}
