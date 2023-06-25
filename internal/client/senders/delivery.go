package senders

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
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
	zipper     mgzip.Zipper
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

func (lm *lineMan) DeliveryBody(mData [][]byte) error {

	var buf bytes.Buffer
	for _, data := range mData {
		buf.Reset()
		request := lm.httpclient.R()
		request.Method = lm.methodSend
		request.URL = lm.receiver
		request.SetHeader("Content-Type", "application/json")

		if lm.zipper != nil {
			v, err := lm.zipper.Compress(data)
			if err != nil {
				logger.Log.Error("can't compress", zap.Error(err))
				continue
			}
			request.SetHeader("Content-Encoding", lm.zipper.GetEncoding())
			buf.Write(v)
		}
		request.SetBody(buf.Bytes())
		response, err := request.Send()
		if err != nil {
			return err
		}
		if response.StatusCode() != http.StatusOK {
			err = fmt.Errorf("%w (%d)!=200", ErrorStatusCode, response.StatusCode())
			logger.Log.Error("DeliveryBody() statusCode: ", zap.Error(err))
			return err
		}

	}
	return nil
}
func NewLineMan(opts ...Option) (DeliveryMan, error) {
	l := &lineMan{
		httpclient: resty.New(),
		zipper:     mgzip.NewGZipper(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}
