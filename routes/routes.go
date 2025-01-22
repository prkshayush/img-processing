package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prkshayush/img-processing/controllers"
)

func ApiRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/submit", controllers.SubmitJob)
		api.GET("/status", controllers.GetJobStatus)
	}
}