package services

import (
	"time"

	"github.com/prkshayush/img-processing/models"
	"github.com/prkshayush/img-processing/rabbitmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	Count  int `json:"count"`
	Visits []struct {
		StoreID   string   `json:"store_id"`
		ImageURLs []string `json:"image_url"`
		VisitTime string   `json:"visit_time"`
	} `json:"visits"`
}

// POST route service
func HandleJobSubmit(request Request) (interface{}, error) {
	job := models.Job{
		Status:    "being processed",
		CreatedAt: time.Now(),
	}

	insertRes, err := models.InsertJobs(rabbitmq.DBCollection, job)
	if err != nil {
		return nil, err
	}

	jobID := insertRes.InsertedID.(primitive.ObjectID).Hex()
    err = rabbitmq.PublishToQueue("jobQueue", jobID)
    if err != nil {
        return nil, err
    }

    return jobID, nil
}

// GET route service
func GetJobStatus(jobID interface{}) (string, error) {
	job, err := models.GetStatusByID(rabbitmq.DBCollection, jobID)
	if err != nil {
		return "", err
	}

	return job.Status, nil
}
