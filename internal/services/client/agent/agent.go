package agent

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/grabbers"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/senders"
	"net/http"
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
	now            func() time.Time
	name           string
	metricsState   map[string]string
	metricsNumbers int
	delivery       senders.DeliveryMan
	grabber        grabbers.Grabber
}
type Option func(agent *agentBond)

func SetMetricsNumber(num int) Option {
	return func(agent *agentBond) {
		agent.metricsNumbers = num

	}
}

func SetName(name string) Option {
	return func(agent *agentBond) {
		agent.name = name

	}
}
func InitDeliveryAddress(address string) Option {
	return func(agent *agentBond) {
		delivery, err := senders.NewLineMan(address) //todo: чо-то с ошибкой делать
		if err != nil {
			panic(fmt.Errorf("InitDeliveryAddress:%w", err))
		}
		agent.delivery = delivery

	}
}
func SetReportInterval(sec int64) Option {
	return func(agent *agentBond) {
		agent.reportInterval = time.Duration(sec) * time.Second
	}
}
func SetPollInterval(sec int64) Option {
	return func(agent *agentBond) {
		agent.pollInterval = time.Duration(sec) * time.Second
	}
}

func SetGrabber() Option {
	return func(agent *agentBond) {
		agent.grabber = grabbers.NewRacoon()
	}
}

func SetFunctionGetTime(fc func() time.Time) Option {
	return func(agent *agentBond) {
		agent.now = fc
	}
}
func NewAgentMetric(opts ...Option) Agent {
	const (
		defaultName           = "bond"
		defaultPollInterval   = 2
		defaultReportInterval = 10
		defaultMetricsNumbers = 1
	)
	agent := &agentBond{
		name:           defaultName,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
		grabber:        grabbers.NewRacoon(grabbers.SetAllowMetrics(client.AllowGaugeMetric)),
		metricsState:   make(map[string]string),
		metricsNumbers: defaultMetricsNumbers,
		now:            time.Now,
	}
	for _, opt := range opts {
		opt(agent)
	}
	return agent
}

func (agent *agentBond) Run() {
	startTime := agent.now()
	isAfter := true /*do it for make test stop*/
	for ; isAfter; isAfter = agent.now().After(startTime) {
		snap := agent.grabber.GetSnapshot()
		agent.updateState(snap)
		if agent.now().Sub(startTime) > agent.reportInterval {

			err := agent.delivery.Delivery(agent.metricsState)
			agent.resetState()
			startTime = agent.now()
			if err != nil {
				fmt.Println(err)
			}

		}
		time.Sleep(agent.pollInterval)

	}
}

func (agent *agentBond) resetState() {
	agent.metricsState = make(map[string]string, agent.metricsNumbers)
}
func (agent *agentBond) updateState(state map[string]string) {
	for key, val := range state {
		agent.metricsState[key] = val
	}
}
