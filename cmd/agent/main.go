package main

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/config/agent"
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
)

func main() {
	agentConfig := config.NewAgentConfig()
	config.ParseFlag(agentConfig)
	pollInterval := agentConfig.GetPollIntervalSecond()
	reportIterval := agentConfig.GetReportIntervalSecond()
	fmt.Println("config agent", agentConfig)
	agent := client.NewAgentMetric(
		client.SetName("anton"),
		client.InitDeliveryAddress(agentConfig.GetMyServer()),
		client.SetReportInterval(reportIterval),
		client.SetPollInterval(pollInterval),
	)
	agent.Run()
}
