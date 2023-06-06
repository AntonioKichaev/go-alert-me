package client

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
}

type LineMan struct {
	receiver   string
	httpclient *resty.Client
}

var ErrorStatusCode = errors.New("delivery status code")

func (lm *LineMan) Delivery(data map[string]string) error {
	for metricType, value := range data {
		urlPath, err := url.JoinPath(lm.receiver, metricType, value)
		if err != nil {
			return err
		}
		request := lm.httpclient.R()
		request.Method = _methodRequestSend
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
func NewLineMan(receiver string) (DeliveryMan, error) {
	return &LineMan{
		receiver:   receiver,
		httpclient: resty.New(),
	}, nil
}
