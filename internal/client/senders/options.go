package senders

import (
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/go-resty/resty/v2"
)

type Option func(l *lineMan)

func SetReceiver(receiver string) Option {
	return func(l *lineMan) {
		l.receiver = receiver
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
