package senders

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
)

//go:generate mockery  --name DeliveryMan
type DeliveryMan interface {
	Delivery(map[string]string) error
	DeliveryBody([][]byte) error
}

type lineMan struct {
	receiver   string
	httpclient *resty.Client
	methodSend string
}

var ErrorStatusCode = errors.New("delivery status code")

func (lm *lineMan) Delivery(data map[string]string) error {
	for metricType, value := range data {
		urlPath, err := url.JoinPath(lm.receiver, metricType, value)
		if err != nil {
			return err
		}
		request := lm.httpclient.R()
		request.Method = lm.methodSend
		request.URL = urlPath
		response, err := request.Send()
		if err != nil {
			return err
		}
		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("%w (%d)!=200", ErrorStatusCode, response.StatusCode())
		}

	}
	return nil
}
func (lm *lineMan) DeliveryBody(data [][]byte) error {
	for _, value := range data {
		request := lm.httpclient.R()
		request.Method = lm.methodSend
		request.URL = lm.receiver
		request.SetHeader("Content-Type", "application/json")
		request.SetBody(value)
		response, err := request.Send()
		if err != nil {
			return err
		}
		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("%w (%d)!=200", ErrorStatusCode, response.StatusCode())
		}

	}
	return nil
}
func NewLineMan(receiver, method string) (DeliveryMan, error) {
	return &lineMan{
		receiver:   receiver,
		httpclient: resty.New(),
		methodSend: method,
	}, nil
}
