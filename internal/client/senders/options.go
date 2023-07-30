package senders

import (
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Option func(l *lineMan)

func SetEndpointRaw(endpoint string) Option {
	return func(l *lineMan) {
		l.endpointRawData = endpoint
	}
}

func SetEndpointJSONData(endpoint string) Option {
	return func(l *lineMan) {
		l.endpointJSONData = endpoint
	}
}

func SeHTTPClient(client *resty.Client) Option {
	return func(l *lineMan) {
		l.httpclient = client
	}
}

func SetMethodSend(methodSend string) Option {
	return func(l *lineMan) {
		l.methodSend = methodSend
	}
}

func SetZipper(zipper mgzip.Zipper) Option {
	return func(l *lineMan) {
		l.zipper = zipper
	}
}

func SetLogger(log *zap.Logger) Option {
	return func(l *lineMan) {
		l.logger = log
	}
}
