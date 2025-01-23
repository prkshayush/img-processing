package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prkshayush/img-processing/services"
)

// validate request & send to rabit queue
func SubmitJob(c *gin.Context) {
	var request services.Request

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Job",
		})
		
		return
	}

	jobID, err := services.HandleJobSubmit(request)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit job",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job_id": jobID,
	})

}

// get request handler
func GetJobStatus(c *gin.Context) {
	jobID := c.DefaultQuery("jobID", "")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing job ID",
		})
	}

	status, err := services.GetJobStatus(jobID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get job stauts",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}