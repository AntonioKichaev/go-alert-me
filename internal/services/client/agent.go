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

func NewAgentMetric(name string, serverAddr string, reportInterval, pollInterval time.Duration) (Agent, error) {
	delivery, err := NewLineMan(serverAddr)
	if err != nil {
		return nil, err
	}
	return &agentBond{
		name:           name,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		grabber:        NewRacoon(),
		delivery:       delivery,
		metricsState:   make(map[string]string, _metricsLenght),
		now:            time.Now,
	}, nil
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
