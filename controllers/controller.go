package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/prkshayush/img-processing/models"
    "github.com/prkshayush/img-processing/services"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// submit controller
func SubmitJob(c *gin.Context) {
    var jobRequest models.JobRequest
    if err := c.ShouldBindJSON(&jobRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    if jobRequest.Count != len(jobRequest.Visits) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Count does not match the number of visits"})
        return
    }

    for _, visit := range jobRequest.Visits {
        if visit.StoreID == "" || visit.VisitTime == "" || len(visit.ImageURLs) == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visit details"})
            return
        }
    }

    jobID, err := services.HandleJobSubmit(jobRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"job_id": jobID})
}

// status controller
func GetJobStatus(c *gin.Context) {
    jobID := c.Query("jobID")
    if jobID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job ID"})
        return
    }

    jobObjectID, err := primitive.ObjectIDFromHex(jobID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
        return
    }

    job, err := models.GetJobByID(jobObjectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": job.Status,
        "job_id": jobID,
        "failed_id": job.FailedID,
    })
}