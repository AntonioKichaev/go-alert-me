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
type Option func(agent *agentBond) error

func SetName(name string) Option {
	return func(agent *agentBond) error {
		agent.name = name
		return nil
	}
}
func InitDeliveryAddress(address string) Option {
	return func(agent *agentBond) error {
		delivery, err := NewLineMan(address) //todo: чо-то с ошибкой делать
		agent.delivery = delivery
		return err
	}
}
func SetReportInterval(sec int64) Option {
	return func(agent *agentBond) error {
		agent.reportInterval = time.Duration(sec) * time.Second
		return nil
	}
}
func SetPollInterval(sec int64) Option {
	return func(agent *agentBond) error {
		agent.pollInterval = time.Duration(sec) * time.Second
		return nil
	}
}

func SetGrabber() Option {
	return func(agent *agentBond) error {
		agent.grabber = NewRacoon()
		return nil
	}
}
func SetMetricState() Option {
	return func(agent *agentBond) error {
		agent.metricsState = make(map[string]string, _metricsLenght)
		return nil
	}
}

func SetFunctionGetTime(fc func() time.Time) Option {
	return func(agent *agentBond) error {
		agent.now = fc
		return nil
	}
}
func NewAgentMetric(opts ...Option) (Agent, error) {
	agent := &agentBond{}
	err := SetName("rand")(agent)
	err = SetPollInterval(2)(agent)
	err = SetReportInterval(10)(agent)
	err = InitDeliveryAddress("localhost:8080")(agent)
	err = SetGrabber()(agent)
	err = SetMetricState()(agent)
	err = SetFunctionGetTime(time.Now)(agent)
	if err != nil {
		panic(fmt.Errorf("newAgentMetric: %v", err))
	}
	for _, opt := range opts {
		err = opt(agent)
		if err != nil {
			panic(fmt.Errorf("newAgentMetric: %v", err))
		}
	}
	return agent, nil
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
