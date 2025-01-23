package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Job struct {
	JobID     interface{} `json:"job_id" bson:"_id,omitempty"`
	Status    string      `json:"status" bson:"status"`
	StoreID   []string    `json:"store_id,omitempty" bson:"store_id,omitempty"`
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

func UpdateStatus(collection *mongo.Collection, jobID interface{}, status string, storeID []string, failedID []string) error {
	filter := bson.M{"_id": jobID}
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"store_id":  storeID,
			"failed_id": failedID,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}
