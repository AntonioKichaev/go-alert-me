package client

import (
	"fmt"
	"net/http"
	"time"
)

const (
	_methodRequestSend = http.MethodPost
	_reciver           = "http://0.0.0.0:8080/update"
)

var _metricsLenght = len(_allowGaugeMetric)

//go:generate mockery --name Agent
type Agent interface {
	Run()
}

type agentBond struct {
	pollIntervalMillis   time.Duration
	reportIntervalMillis time.Duration
	now                  func() time.Time
	name                 string
	metricsState         map[string]string
	delivery             DeliveryMan
	grabber              Grabber
}

func NewAgentMetric(name string, reportIntervalMillis, pollIntervalMillis time.Duration) Agent {
	return &agentBond{
		name:                 name,
		pollIntervalMillis:   pollIntervalMillis,
		reportIntervalMillis: reportIntervalMillis,
		grabber:              NewRacoon(),
		delivery:             NewLineMan(_reciver),
		metricsState:         make(map[string]string, _metricsLenght),
		now:                  time.Now,
	}
}

func (agent *agentBond) Run() {
	startTime := agent.now()
	isAfter := true /*do it for make test stop*/
	for ; isAfter; isAfter = agent.now().After(startTime) {
		snap := agent.grabber.GetSnapshot()
		agent.updateState(snap)
		if agent.now().Sub(startTime) > agent.reportIntervalMillis {

			err := agent.delivery.Delivery(agent.metricsState)
			agent.resetState()
			startTime = agent.now()
			if err != nil {
				fmt.Println(err)
			}

		}
		time.Sleep(agent.pollIntervalMillis)

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
