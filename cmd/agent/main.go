package main

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/config/agent"
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
	"net/url"
)

const _endPointUpdateValue = "/update"

func main() {
	agentConfig := config.NewAgentConfig()
	config.ParseFlag(agentConfig)
	pollInterval := agentConfig.GetPollIntervalSecond()
	reportIterval := agentConfig.GetReportIntervalSecond()
	fmt.Println("config agent", agentConfig)
	diliveryAddress, err := url.JoinPath(agentConfig.GetMyServer(), _endPointUpdateValue)
	if err != nil {
		panic(err)
	}
	agent := client.NewAgentMetric(
		client.SetName("anton"),
		client.InitDeliveryAddress(diliveryAddress),
		client.SetReportInterval(reportIterval),
		client.SetPollInterval(pollInterval),
	)
	agent.Run()
}
