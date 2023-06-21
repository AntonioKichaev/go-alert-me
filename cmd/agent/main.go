package main

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/config/agent"
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/agent"
	"net/url"
)

const _endPointUpdateValue = "/update"

func main() {
	agentConfig := config.NewAgentConfig()
	config.ParseFlag(agentConfig)
	pollInterval := agentConfig.GetPollIntervalSecond()
	reportInterval := agentConfig.GetReportIntervalSecond()
	fmt.Println("config agent", agentConfig)
	diliveryAddress, err := url.JoinPath(agentConfig.GetMyServer(), _endPointUpdateValue)
	if err != nil {
		panic(err)
	}
	ag := agent.NewAgentMetric(
		agent.SetName("anton"),
		agent.InitDeliveryAddress(diliveryAddress),
		agent.SetReportInterval(reportInterval),
		agent.SetPollInterval(pollInterval),
		agent.SetMetricsNumber(len(client.AllowGaugeMetric)),
	)
	ag.Run()
}
