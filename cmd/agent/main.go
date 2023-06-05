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
	fmt.Println("config ", agentConfig)
	agent, err := client.NewAgentMetric(
		"",
		agentConfig.GetMyServer(),
		reportIterval,
		pollInterval,
	)
	if err != nil {
		panic(err)
	}
	agent.Run()
}
