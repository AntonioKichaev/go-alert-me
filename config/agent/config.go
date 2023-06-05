package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Agent struct {
	httpServerAdr        string
	reportIntervalSecond int64
	pollIntervalSecond   int64
}

func (a *Agent) GetReportIntervalSecond() time.Duration {
	return time.Duration(a.reportIntervalSecond) * time.Second
}
func (a *Agent) GetPollIntervalSecond() time.Duration {
	return time.Duration(a.pollIntervalSecond) * time.Second
}
func (a *Agent) GetMyServer() string {
	if strings.Contains(a.httpServerAdr, "localhost") {
		result, _ := url.JoinPath("http://", a.httpServerAdr)
		return result
	}
	return a.httpServerAdr
}
func (a *Agent) String() string {
	return fmt.Sprintf("server:(%s)\nreportInterval:(%d sec)\npollInterval:(%d sec)", a.httpServerAdr, a.reportIntervalSecond, a.pollIntervalSecond)
}

func NewAgentConfig() *Agent {
	return &Agent{}
}
