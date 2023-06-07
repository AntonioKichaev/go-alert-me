package client

import (
	"fmt"
	"net/http"
	"time"
)

const (
	_methodRequestSend = http.MethodPost
)

var _metricsLenght = len(_allowGaugeMetric)

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
	delivery       DeliveryMan
	grabber        Grabber
}
type Option func(agent *agentBond)

func SetName(name string) Option {
	return func(agent *agentBond) {
		agent.name = name

	}
}
func InitDeliveryAddress(address string) Option {
	return func(agent *agentBond) {
		delivery, err := NewLineMan(address) //todo: чо-то с ошибкой делать
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
		agent.grabber = NewRacoon()
	}
}
func SetMetricState() Option {
	return func(agent *agentBond) {
		agent.metricsState = make(map[string]string, _metricsLenght)
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
	)
	agent := &agentBond{
		name:           defaultName,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
		grabber:        NewRacoon(),
		metricsState:   make(map[string]string, _metricsLenght),
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
	agent.metricsState = make(map[string]string, _metricsLenght)
}
func (agent *agentBond) updateState(state map[string]string) {
	for key, val := range state {
		agent.metricsState[key] = val
	}
}
