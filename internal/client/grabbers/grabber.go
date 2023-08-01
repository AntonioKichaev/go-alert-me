package grabbers

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

const _additionsMetrics = 5

//go:generate mockery  --name Grabber
type Grabber interface {
	GetSnapshot() map[string]string
}

//go:generate mockery  --name Random
type Random interface {
	Int() int
}
type racoon struct {
	random           Random
	allowMetrics     map[string]struct{}
	isolationMetrics map[string]string
	mu               sync.Mutex
}

func NewRacoon(opts ...Option) Grabber {
	rn := &racoon{
		random:           rand.New(rand.NewSource(322)),
		mu:               sync.Mutex{},
		isolationMetrics: make(map[string]string, _additionsMetrics),
	}
	for _, opt := range opts {
		opt(rn)
	}

	go rn.setAdditionalMetricsGoUtil()
	return rn
}

func (rc *racoon) GetSnapshot() map[string]string {

	snapshot := make(map[string]string, len(rc.allowMetrics)+_additionsMetrics)
	rc.setGauge(snapshot)
	rc.setAdditionalMetrics(snapshot)
	rc.mu.Lock()
	defer rc.mu.Unlock()
	for key, val := range rc.isolationMetrics {
		snapshot[key] = val
	}
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

func (rc *racoon) setAdditionalMetricsGoUtil() {
	v, _ := mem.VirtualMemory()

	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.isolationMetrics["gauge/Total"] = strconv.FormatUint(v.Total, 10)
	rc.isolationMetrics["gauge/FreeMemory"] = strconv.FormatUint(v.Free, 10)
	c, _ := cpu.Counts(true)
	rc.isolationMetrics["gauge/CPUutilization1"] = strconv.Itoa(c)
}
