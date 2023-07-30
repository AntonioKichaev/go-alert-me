package senders

import (
	metricsEntity "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/pkg/hasher"
)

type Hasher interface {
	Sign(data []byte) string
}

type lineManHasher struct {
	DeliveryMan
	key string
	h   Hasher
}

func NewLineManHasher(dm DeliveryMan, key string) (*lineManHasher, error) {
	h := hasher.NewHasher(key)
	lmh := &lineManHasher{
		DeliveryMan: dm,
		key:         key,
		h:           h,
	}
	return lmh, nil
}

func (lm *lineManHasher) Delivery(m map[string]string) error {
	return lm.DeliveryMan.Delivery(m)
}
func (lm *lineManHasher) DeliveryBody(data [][]byte) error {
	return lm.DeliveryMan.DeliveryBody(data)
}
func (lm *lineManHasher) DeliveryMetricsJSON(m []metricsEntity.Metrics) error {
	return lm.DeliveryMan.DeliveryMetricsJSON(m)
}
