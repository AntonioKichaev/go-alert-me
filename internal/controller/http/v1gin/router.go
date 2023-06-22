package v1gin

import (
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.Engine, up usecase.Updater, metric usecase.ReceiverMetric) {
	rg := route.Group("/v1")
	NewUpdaterRoutes(rg, up)
}
