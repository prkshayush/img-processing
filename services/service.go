package services

import (
	"time"

	"github.com/prkshayush/img-processing/models"
	"github.com/prkshayush/img-processing/rabbitmq"
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

	err = rabbitmq.PublishToQueue("jobQueue", insertRes.InsertedID)
	if err != nil {
		return nil, err
	}

	return insertRes.InsertedID, nil
}

// GET route service
func GetJobStatus(jobID interface{}) (string, error) {
	job, err := models.GetStatusByID(rabbitmq.DBCollection, jobID)
	if err != nil {
		return "", err
	}

	return job.Status, nil
}
