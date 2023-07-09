package agent

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/client/grabbers"
	"github.com/antoniokichaev/go-alert-me/internal/client/senders"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"go.uber.org/zap"
	"time"
)

type Option func(agent *agentBond)

func WithLogger(logger *zap.Logger) Option {
	return func(agent *agentBond) {
		agent.logger = logger
	}
}
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
func InitDeliveryAddress(endpointRawData, endpointJSONData, method string) Option {
	return func(agent *agentBond) {
		delivery, err := senders.NewLineMan(
			senders.SetEndpointJSONData(endpointJSONData),
			senders.SetEndpointRaw(endpointRawData),
			senders.SetMethodSend(method),
			senders.SetZipper(agent.zipper),
			senders.SetLogger(agent.logger),
		) //todo: чо-то с ошибкой делать
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
func SetZipper(zipper mgzip.Zipper) Option {
	return func(agent *agentBond) {
		agent.zipper = zipper

	}
}
