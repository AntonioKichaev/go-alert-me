package v1gin

import (
	"errors"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
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

type UpdaterRoutes struct {
	uc usecase.Updater
}

func NewUpdaterRoutes(uc usecase.Updater) *UpdaterRoutes {
	return &UpdaterRoutes{uc: uc}

}

// updateMetrics принимает запрос ввида /update/{counter|gauge}/someMetric/527
func (h *UpdaterRoutes) updateMetrics(ctx *gin.Context) {
	metricType := ctx.Param(_metricType)
	metricName := ctx.Param(_metricName)
	metricValue := ctx.Param(_metricValue)

	_, err := h.uc.UpdateMetricByParams(metricName, metricType, metricValue)
	if err != nil {
		if errors.Is(err, metrics2.ErrorName) {
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
