package main

import (
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
	"time"
)

func main() {
	pollInterval := time.Second * 2
	reportIterval := time.Second * 10

	agent := client.NewAgentMetric(
		"",
		reportIterval,
		pollInterval,
	)
	agent.Run()
}
