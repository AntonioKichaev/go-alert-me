package agent

import (
	"encoding/json"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/grabbers"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/senders"
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
	"net/http"
	"strconv"
	"strings"
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
}
type Option func(agent *agentBond)

func SetNotifyChan(ch chan struct{}) Option {
	return func(agent *agentBond) {
		agent.notify = ch
	}
}
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
func InitDeliveryAddress(address, method string) Option {
	return func(agent *agentBond) {
		delivery, err := senders.NewLineMan(address, method) //todo: чо-то с ошибкой делать
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
	reportTicker := time.NewTicker(agent.reportInterval)
	pollTicker := time.NewTicker(agent.pollInterval)
	for {
		select {
		case <-reportTicker.C:
			data := agent.makeFormatToSend()
			if len(data) > 0 {
				err := agent.delivery.DeliveryBody(data)
				fmt.Println(err)
			} else {

				fmt.Println("Nothind to send")
			}
			agent.resetState()
		case <-pollTicker.C:
			snap := agent.grabber.GetSnapshot()
			agent.updateState(snap)
		case <-agent.notify:
			return
		default:
			time.Sleep(time.Second / 2)
		}
	}
}

func (agent *agentBond) resetState() {
	agent.metricsState = make(map[string]string, agent.metricsNumbers)
}
func (agent *agentBond) updateState(state map[string]string) {
	for key, val := range state {
		if strings.Contains(key, "counter") {
			nVal, _ := strconv.Atoi(val)
			oldVal, _ := strconv.Atoi(agent.metricsState[key])
			val = strconv.Itoa(nVal + oldVal)
		}
		agent.metricsState[key] = val
	}
}

func (agent *agentBond) makeFormatToSend() [][]byte {
	// 1 way
	// передавать строки ввида json {}
	// надо будет чтоб агент превращал state ->json
	//
	//
	res := make([][]byte, 0, len(agent.metricsState))

	for key, val := range agent.metricsState {
		s := strings.Split(key, "/")
		if len(s) < 2 {
			continue
		}
		t, err := metrics.NewMetrics(s[0], s[1], val)
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
