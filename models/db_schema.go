package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Job struct {
	JobID     int       `json:"job_id" bson:"job_id"`
	Status    string    `json:"status" bson:"status"`
	StoreID   []string  `json:"store_id,omitempty" bson:"store_id,omitempty"`
	FailedID  []string  `json:"failed_id,omitempty" bson:"failed_id,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// POST insert function
func InsertJobArr(collection *mongo.Collection, job Job) error{
	_, err := collection.InsertOne(context.TODO(), job)

	return err
}

// GET by Id function
func GetStatusByID(collection *mongo.Collection, jobID Job) (*Job, error) {
	var job Job
	err := collection.FindOne(context.TODO(), bson.M{"job_id": jobID}).Decode(&job)

	return &job, err
}

// PATCH functions

func UpdateStatus(collection *mongo.Collection, jobID Job, status string) error {
	filter := bson.M{"job_id": jobID}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func AddStoreID(collection *mongo.Collection, jobID Job, storeID string) error{
	filter := bson.M{"job_id": jobID}
	update := bson.M{"$addToSet": bson.M{"store_ids": storeID}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func AddFailedID(collection *mongo.Collection, jobID Job, failedID string) error{
	filter := bson.M{"job_id": jobID}
	update := bson.M{"$addToSet": bson.M{"failed_id": failedID}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	return err
}