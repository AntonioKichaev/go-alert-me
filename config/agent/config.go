package config

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"net/url"
	"strings"
)

type Agent struct {
	HTTPServerAdr        string `env:"ADDRESS"`
	ReportIntervalSecond int64  `env:"REPORT_INTERVAL"`
	PollIntervalSecond   int64  `env:"POLL_INTERVAL"`
	LoggingLevel         string `env:"LOGGING_LEVEL"`
	SecretKey            string `env:"KEY"`
}

func (a *Agent) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("HTTPServerAdr", a.HTTPServerAdr)
	encoder.AddInt64("ReportIntervalSecond", a.ReportIntervalSecond)
	encoder.AddInt64("PollIntervalSecond", a.PollIntervalSecond)
	encoder.AddString("LOGGING_LEVEL", a.LoggingLevel)
	encoder.AddString("SecretKey", a.SecretKey)
	return nil
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
func (a *Agent) GetLoggingLevel() string {
	return a.LoggingLevel
}

func (a *Agent) String() string {
	return fmt.Sprintf("server:(%s)\nreportInterval:(%d sec)\npollInterval:(%d sec)", a.HTTPServerAdr, a.ReportIntervalSecond, a.PollIntervalSecond)
}

func NewAgentConfig(opts ...Option) *Agent {
	const (
		_defaultHTTPServerAdr        = "localhost:8080"
		_defaultReportIntervalSecond = 10
		_defaultPollIntervalSecond   = 2
		_defaultLoggingLevel         = "INFO"
	)
	agent := &Agent{
		HTTPServerAdr:        _defaultHTTPServerAdr,
		ReportIntervalSecond: _defaultReportIntervalSecond,
		PollIntervalSecond:   _defaultPollIntervalSecond,
		LoggingLevel:         _defaultLoggingLevel,
	}

	for _, opt := range opts {
		opt(agent)
	}
	return agent
}
