package v1gin

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.Engine, up usecase.Updater, metric usecase.ReceiverMetric) {
	rg := route.Group("/v1")
	ur := NewUpdaterRoutes(up)
	g := rg.Group("/update")
	{
		g.POST(fmt.Sprintf("/:%s/:%s/:%s", _metricType, _metricName, _metricValue), ur.updateMetrics)
	}

}
