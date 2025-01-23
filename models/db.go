package models

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DBCollection *mongo.Collection

func ConnectDB(uri, dbName, collectionName string) error {
    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return err
    }

    DBCollection = client.Database(dbName).Collection(collectionName)
    return nil
}

func InsertJob(job Job) (*mongo.InsertOneResult, error) {
    return DBCollection.InsertOne(context.TODO(), job)
}

func UpdateJobStatus(jobID primitive.ObjectID, status string, failedStoreID []string) error {
    update := bson.M{"$set": bson.M{"status": status}}
    if len(failedStoreID) > 0 {
        update["$set"].(bson.M)["failed_id"] = failedStoreID
    }
    _, err := DBCollection.UpdateOne(context.TODO(), bson.M{"_id": jobID}, update)
    return err
}

func GetJobByID(jobID primitive.ObjectID) (*Job, error) {
    var job Job
    err := DBCollection.FindOne(context.TODO(), bson.M{"_id": jobID}).Decode(&job)
    return &job, err
}