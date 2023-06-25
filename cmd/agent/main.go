package main

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/config/agent"
	"github.com/antoniokichaev/go-alert-me/internal/client"
	"github.com/antoniokichaev/go-alert-me/internal/client/agent"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"net/http"
	"net/url"
)

const _endPointUpdateValue = "/update/"

func main() {
	agentConfig := config.NewAgentConfig()
	config.ParseFlag(agentConfig)
	pollInterval := agentConfig.GetPollIntervalSecond()
	reportInterval := agentConfig.GetReportIntervalSecond()
	fmt.Println("config agent", agentConfig)
	deliveryAddress, err := url.JoinPath(agentConfig.GetMyServer(), _endPointUpdateValue)
	zipper := mgzip.NewGZipper()
	if err != nil {
		panic(err)
	}
	ag := agent.NewAgentMetric(
		agent.SetName("anton"),
		agent.SetZipper(zipper),
		agent.InitDeliveryAddress(deliveryAddress, http.MethodPost),
		agent.SetReportInterval(reportInterval),
		agent.SetPollInterval(pollInterval),
		agent.SetMetricsNumber(len(client.AllowGaugeMetric)),
	)
	ag.Run()
}
