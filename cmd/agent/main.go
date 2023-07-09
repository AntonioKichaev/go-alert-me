package main

import (
	"github.com/antoniokichaev/go-alert-me/config/agent"
	"github.com/antoniokichaev/go-alert-me/internal/client"
	"github.com/antoniokichaev/go-alert-me/internal/client/agent"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

const (
	_endPointUpdateValue  = "/update/"
	_endPointUpdateValues = "/updates/"
)

func main() {
	agentConfig := config.NewAgentConfig()
	config.ParseFlag(agentConfig)
	l := logger.Initialize(agentConfig.LoggingLevel)
	pollInterval := agentConfig.GetPollIntervalSecond()
	reportInterval := agentConfig.GetReportIntervalSecond()
	l.Info("config agent", zap.Object("agent", agentConfig))
	deliveryAddress, err := url.JoinPath(agentConfig.GetMyServer(), _endPointUpdateValue)
	if err != nil {
		panic(err)
	}
	deliveryAddressJSON, err := url.JoinPath(agentConfig.GetMyServer(), _endPointUpdateValues)
	zipper := mgzip.NewGZipper()
	if err != nil {
		panic(err)
	}
	ag := agent.NewAgentMetric(
		agent.WithLogger(l),
		agent.SetName("anton"),
		agent.SetZipper(zipper),
		agent.InitDeliveryAddress(deliveryAddress, deliveryAddressJSON, http.MethodPost),
		agent.SetReportInterval(reportInterval),
		agent.SetPollInterval(pollInterval),
		agent.SetMetricsNumber(len(client.AllowGaugeMetric)),
	)
	ag.Run()
}
