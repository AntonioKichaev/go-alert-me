package v1gin

import (
	"errors"
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/entity"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	//При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	_zeroMetricName = http.StatusNotFound
	//При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
	_incorrectMetricType         = http.StatusBadRequest
	_incorrectMetricValue        = http.StatusBadRequest
	_contentTypeText             = "text/plain; charset=utf-8"
	_metricType           string = "MetricType"
	_metricName                  = "MetricName"
	_metricValue                 = "MetricValue"
)

type updaterRoutes struct {
	uc usecase.Updater
}

func NewUpdaterRoutes(handler gin.IRouter, uc usecase.Updater) {
	ur := newUpdaterRoutes(uc)

	handler.POST(fmt.Sprintf("/update/:%s/:%s/:%s", _metricType, _metricName, _metricValue), ur.updateMetrics)
}
func newUpdaterRoutes(uc usecase.Updater) *updaterRoutes {
	return &updaterRoutes{uc: uc}
}

// updateMetrics принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *updaterRoutes) updateMetrics(ctx *gin.Context) {
	metricType := ctx.Param(_metricType)
	metricName := ctx.Param(_metricName)
	metricValue := ctx.Param(_metricValue)

	var err error

	switch entity.MetricType(metricType) {
	case entity.GaugeName:
		err = h.uc.SetGauge(metricName, metricValue)
	case entity.CounterName:
		err = h.uc.AddCounter(metricName, metricValue)
	default:
		err = ErrorUnknownMetricType
	}
	if err != nil {
		if errors.Is(err, entity.ErrorName) {
			ctx.AbortWithStatus(_zeroMetricName)
			return
		}
		if errors.Is(err, ErrorUnknownMetricType) {
			ctx.AbortWithStatus(_incorrectMetricType)
			return
		}
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.Data(http.StatusOK, "text/plain; charset=utf-8", nil)

}
