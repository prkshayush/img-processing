package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageResult struct {
	ImageURL  string `json:"image_url" bson:"image_url"`
	Perimeter int    `json:"perimeter" bson:"perimeter"`
	Processed bool   `json:"processed" bson:"processed"`
}

type Visit struct {
	StoreID   string        `json:"store_id" bson:"store_id"`
	ImageURLs []string      `json:"image_url" bson:"image_url"`
	VisitTime string        `json:"visit_time" bson:"visit_time"`
	Results   []ImageResult `json:"results,omitempty" bson:"results,omitempty"`
}

type Job struct {
	JobID     primitive.ObjectID `json:"job_id" bson:"_id,omitempty"`
	Status    string             `json:"status" bson:"status"`
	Visits    []Visit            `json:"visits,omitempty" bson:"visits,omitempty"`
	FailedID  []string           `json:"failed_id,omitempty" bson:"failed_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type JobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}
