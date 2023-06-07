package config

import (
	"fmt"
	"net/url"
	"strings"
)

type Agent struct {
	HTTPServerAdr        string `env:"ADDRESS"`
	ReportIntervalSecond int64  `env:"REPORT_INTERVAL"`
	PollIntervalSecond   int64  `env:"POLL_INTERVAL"`
}

func (a *Agent) GetReportIntervalSecond() int64 {
	return a.ReportIntervalSecond
}
func (a *Agent) GetPollIntervalSecond() int64 {
	return a.PollIntervalSecond
}
func (a *Agent) GetMyServer() string {
	if strings.Contains(a.HTTPServerAdr, "localhost") {
		result, _ := url.JoinPath("http://", a.HTTPServerAdr)
		return result
	}
	return a.HTTPServerAdr
}
func (a *Agent) String() string {
	return fmt.Sprintf("server:(%s)\nreportInterval:(%d sec)\npollInterval:(%d sec)", a.HTTPServerAdr, a.ReportIntervalSecond, a.PollIntervalSecond)
}

func NewAgentConfig() *Agent {
	agent := &Agent{}

	return agent
}
