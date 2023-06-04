package main

import (
	"github.com/antoniokichaev/go-alert-me/internal/services/client"
	"time"
)

func main() {
	pollIntervalMillis := time.Second * 2
	reportItervalMillis := time.Second * 10

	agent := client.NewAgentMetric(
		"",
		reportItervalMillis,
		pollIntervalMillis,
	)
	agent.Run()
}
