package core

import (
	"github.com/gin-gonic/gin"
)

//RegisterCoreRoutes registers core routes with server
func RegisterCoreRoutes(r *gin.Engine, controller *Controller) {
	r.POST("/deals", getDealsValidator, controller.GetFilteredDeals)
}
