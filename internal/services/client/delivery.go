package client

import (
	"net/http"
	"net/url"
	"time"
)

//go:generate mockery  --name DeliveryMan
type DeliveryMan interface {
	Delivery(map[string]string) error
}

type LineMan struct {
	receiver   string
	httpclient http.Client
}

func (lm *LineMan) Delivery(data map[string]string) error {
	for key, val := range data {
		urlPath, err := url.JoinPath(lm.receiver, key, val)
		if err != nil {
			return err
		}
		request, err := http.NewRequest(_methodRequestSend, urlPath, nil)
		if err != nil {
			return err
		}
		resp, err := lm.httpclient.Do(request)

		if err != nil {
			// server isn't available
			return err
		}
		_ = resp.Body.Close()

	}
	return nil
}
func NewLineMan(receiver string) (DeliveryMan, error) {
	return &LineMan{
		receiver:   receiver,
		httpclient: http.Client{Timeout: time.Second * 2},
	}, nil
}
