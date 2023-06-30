package agent

import (
	"encoding/json"
	"github.com/antoniokichaev/go-alert-me/internal/client"
	"github.com/antoniokichaev/go-alert-me/internal/client/grabbers"
	"github.com/antoniokichaev/go-alert-me/internal/client/senders"
	"github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MethodRequestSend = http.MethodPost
)

//go:generate mockery --name Agent
type Agent interface {
	Run()
}

type agentBond struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	name           string
	metricsState   map[string]string
	metricsNumbers int
	delivery       senders.DeliveryMan
	grabber        grabbers.Grabber
	notify         <-chan struct{}
	zipper         mgzip.Zipper
	mu             sync.RWMutex
	logger         *zap.Logger
}

func NewAgentMetric(opts ...Option) Agent {
	const (
		defaultName           = "bond"
		defaultPollInterval   = 2 * time.Second
		defaultReportInterval = 10 * time.Second
		defaultMetricsNumbers = 1
	)
	agent := &agentBond{
		name:           defaultName,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
		grabber:        grabbers.NewRacoon(grabbers.SetAllowMetrics(client.AllowGaugeMetric)),
		metricsState:   make(map[string]string),
		metricsNumbers: defaultMetricsNumbers,
		notify:         make(chan struct{}),
	}
	for _, opt := range opts {
		opt(agent)
	}
	return agent
}

func (agent *agentBond) Run() {
	go agent.sendReport()
	go agent.updateState()
	<-agent.notify // ждем сигнала на завершения

}

func (agent *agentBond) sendReport() {
	for ; ; time.Sleep(agent.reportInterval) {
		data := agent.makeFormatToSend()
		if len(data) > 0 {
			err := agent.delivery.DeliveryBody(data)
			if err != nil {
				agent.logger.Error("agent.Run() delivery err:=", zap.Error(err))
			}
		}
		agent.resetState()
	}

}

func (agent *agentBond) resetState() {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	agent.metricsState = make(map[string]string, agent.metricsNumbers)
}
func (agent *agentBond) updateState() {
	for ; ; time.Sleep(agent.pollInterval) {
		state := agent.grabber.GetSnapshot()
		for key, val := range state {
			if strings.Contains(key, "counter") {
				nVal, err := strconv.Atoi(val)
				if err != nil {
					agent.logger.Info("val: cant convert to integer", zap.String("val", val), zap.Error(err))
					continue
				}
				agent.mu.RLock()
				oldV, ok := agent.metricsState[key]
				agent.mu.RUnlock()
				var oldVal int
				if ok {
					oldVal, err = strconv.Atoi(oldV)
					if err != nil {
						agent.logger.Info("agent.MetricsState: cant convert to integer", zap.String("oldV", oldV), zap.Error(err))
						continue
					}
				}
				val = strconv.Itoa(nVal + oldVal)
			}
			agent.mu.Lock()
			agent.metricsState[key] = val
			agent.mu.Unlock()
		}
	}
}

func (agent *agentBond) makeFormatToSend() [][]byte {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	if len(agent.metricsState) == 0 {
		return nil
	}
	res := make([][]byte, 0, len(agent.metricsState))
	for key, val := range agent.metricsState {
		s := strings.Split(key, "/")
		if len(s) < 2 {
			continue
		}
		t, err := metrics.NewMetrics(
			metrics.SetMetricType(s[0]),
			metrics.SetName(s[1]),
			metrics.SetValueOrDelta(val))
		if err != nil {
			continue
		}
		b, err := json.Marshal(t)
		if err != nil {
			continue
		}
		res = append(res, b)

	}
	return res
}
