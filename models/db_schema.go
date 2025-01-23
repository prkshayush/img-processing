package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Visit struct {
	StoreID   string   `json:"store_id" bson:"store_id"`
	ImageURLs []string `json:"image_url" bson:"image_url"`
	VisitTime string   `json:"visit_time" bson:"visit_time"`
}

type Job struct {
	JobID     interface{} `json:"job_id" bson:"_id,omitempty"`
	Status    string      `json:"status" bson:"status"`
	Visits    []Visit     `json:"visits" bson:"visits"`
	FailedID  []string    `json:"failed_id,omitempty" bson:"failed_id,omitempty"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
}

// POST insert function
func InsertJobs(collection *mongo.Collection, job Job) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(context.TODO(), job)
}

// GET by Id function
func GetStatusByID(collection *mongo.Collection, jobID interface{}) (*Job, error) {
	var job Job
	err := collection.FindOne(context.TODO(), bson.M{"_id": jobID}).Decode(&job)

	return &job, err
}

// PATCH functions

func UpdateStatus(collection *mongo.Collection, jobID interface{}, status string, failedStoreID string) error {
	filter := bson.M{"_id": jobID}
	update := bson.M{"$set": bson.M{"status": status}}
	if status == "failed" {
		update["$push"] = bson.M{"failed_id": failedStoreID}
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}
