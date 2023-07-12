package senders

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	metricsEntity "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"time"
)

const _maxTrySend = 3

//go:generate mockery  --name DeliveryMan
type DeliveryMan interface {
	Delivery(map[string]string) error
	DeliveryBody([][]byte) error
	DeliveryMetricsJSON([]metricsEntity.Metrics) error
}

type lineMan struct {
	endpointRawData  string
	endpointJSONData string
	httpclient       *resty.Client
	methodSend       string
	zipper           mgzip.Zipper
	logger           *zap.Logger
}

var ErrorStatusCode = errors.New("delivery status code")

func (lm *lineMan) Delivery(data map[string]string) error {
	for metricType, value := range data {
		urlPath, err := url.JoinPath(lm.endpointRawData, metricType, value)
		if err != nil {
			return err
		}
		request := lm.httpclient.R()
		request.Method = lm.methodSend
		request.URL = urlPath
		response, err := lm.Send(request)
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
		request.URL = lm.endpointRawData
		request.SetHeader("Content-Type", "application/json")

		if lm.zipper != nil {
			v, err := lm.zipper.Compress(data)
			if err != nil {
				lm.logger.Error("can't compress", zap.Error(err))
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
			lm.logger.Error("DeliveryBody() statusCode: ", zap.Error(err))
			return err
		}

	}
	return nil
}

func (lm *lineMan) DeliveryMetricsJSON(mSlice []metricsEntity.Metrics) error {
	data, err := json.Marshal(mSlice)
	if err != nil {
		lm.logger.Error("can't Marshal", zap.Error(err))
		return err
	}
	request := lm.httpclient.R()
	request.Method = http.MethodPost
	request.URL = lm.endpointJSONData
	request.SetHeader("Content-Type", "application/json")

	if lm.zipper != nil {
		v, err := lm.zipper.Compress(data)
		if err != nil {
			lm.logger.Error("can't compress", zap.Error(err))
			return err
		}
		request.SetHeader("Content-Encoding", lm.zipper.GetEncoding())
		request.SetBody(v)
	}

	response, err := lm.Send(request)
	if err != nil {
		return fmt.Errorf("send err %w", err)
	}
	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("%w (%d)!=200", ErrorStatusCode, response.StatusCode())
		lm.logger.Error("DeliveryMetricsJSON statusCode: ", zap.Error(err))
		return err
	}
	return nil
}
func (lm *lineMan) Send(request *resty.Request) (response *resty.Response, err error) {
	response = &resty.Response{}
	for i := 1; i <= _maxTrySend; i++ {
		response, err = request.Send()
		if err, ok := err.(net.Error); ok && err.Timeout() && i < _maxTrySend {
			time.Sleep(time.Second * time.Duration(i+i-1))
			continue
		}

		if err != nil {
			return
		}
		break
	}
	return
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
