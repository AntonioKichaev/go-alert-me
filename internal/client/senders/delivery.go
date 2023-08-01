package senders

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	metricsEntity "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/pkg/mgzip"
)

const _maxTrySend = 3

//go:generate mockery  --name DeliveryMan
type DeliveryMan interface {
	Delivery(map[string]string) error
	DeliveryBody([][]byte) error
	DeliveryMetricsJSON([]metricsEntity.Metrics) error
}

type Hasher interface {
	Sign(data []byte) string
}

type lineMan struct {
	endpointRawData  string
	endpointJSONData string
	httpclient       *resty.Client
	methodSend       string
	zipper           mgzip.Zipper
	logger           *zap.Logger
	hash             Hasher
	maxWorkerPool    int
	jobs             chan resty.Request
}

func (lm *lineMan) Delivery(data map[string]string) error {
	for metricType, value := range data {
		urlPath, err := url.JoinPath(lm.endpointRawData, metricType, value)
		if err != nil {
			return err
		}
		request := lm.httpclient.R()
		request.Method = lm.methodSend
		request.URL = urlPath
		lm.AddRequest(request)
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
		if lm.hash != nil {
			sign := lm.hash.Sign(buf.Bytes())
			request.SetHeader("HashSHA256", sign)
		}
		request.SetBody(buf.Bytes())

		lm.AddRequest(request)

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

	var buf bytes.Buffer
	if lm.zipper != nil {
		v, err := lm.zipper.Compress(data)
		if err != nil {
			lm.logger.Error("can't compress", zap.Error(err))
			return err
		}
		request.SetHeader("Content-Encoding", lm.zipper.GetEncoding())
		buf.Write(v)
		request.SetBody(buf.Bytes())
	}
	if lm.hash != nil {
		sign := lm.hash.Sign(buf.Bytes())
		request.SetHeader("HashSHA256", sign)
	}

	lm.AddRequest(request)

	return nil
}

func (lm *lineMan) AddRequest(request *resty.Request) {
	lm.jobs <- *request
}
func (lm *lineMan) Send(request *resty.Request) *resty.Response {
	response := &resty.Response{}
	for i := 1; i <= _maxTrySend; i++ {
		response, err := request.Send()
		if err, ok := err.(net.Error); ok && err.Timeout() && i < _maxTrySend {
			time.Sleep(time.Second * time.Duration(i+i-1))
			continue
		}

		if err != nil {
			lm.logger.Error("error send err", zap.Error(err))
			return response
		}
		break
	}
	return response
}

func NewLineMan(opts ...Option) (DeliveryMan, error) {
	l := &lineMan{
		httpclient:    resty.New(),
		zipper:        mgzip.NewGZipper(),
		maxWorkerPool: 3,
	}

	for _, opt := range opts {
		opt(l)
	}
	jobs := l.CreateWorkerPool(l.maxWorkerPool)
	l.jobs = jobs

	return l, nil
}

func (lm *lineMan) worker(jobs <-chan resty.Request) {
	for {
		req := <-jobs
		lm.Send(&req)

	}
}

func (lm *lineMan) CreateWorkerPool(nWorkerCount int) chan resty.Request {
	jobs := make(chan resty.Request, nWorkerCount)

	for i := 0; nWorkerCount > i; i++ {
		lm.logger.Info("Create worker", zap.Int("current worker", i))
		go lm.worker(jobs)
	}

	return jobs
}
