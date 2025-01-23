package workers

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"

    "github.com/prkshayush/img-processing/models"
    "github.com/prkshayush/img-processing/rabbitmq"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessJobs() {
    msgs, err := rabbitmq.Channel.Consume(
        "jobQueue",
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to register: %v", err)
    }

    // making channel
    processChan := make(chan bool)

    // go routine
    go func() {
        for d := range msgs {
            jobID := string(d.Body)
            processJob(jobID)
        }
    }()

    log.Printf("[*] Waiting for messages, CTRL+C to terminate process")
    // reading from channel
    <-processChan
}

// helper functions

func processJob(jobID string) {
    jobObjectID, _ := primitive.ObjectIDFromHex(jobID)
    job, err := models.GetStatusByID(rabbitmq.DBCollection, jobObjectID)
    if err != nil {
        log.Printf("Failed to get job: %v", err)
        return
    }

    status := "complete"
    for _, visit := range job.Visits {
        for _, imageURL := range visit.ImageURLs {
            err := processImage(imageURL)
            if err != nil {
                status = "failed"
                updateJobStatus(jobID, status, visit.StoreID)
                return
            }
        }
    }

    // status update after processing
    updateJobStatus(jobID, status, "")
}

func processImage(imgURL string) error {
    response, err := http.Get(imgURL)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    // processing job
    time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)

    // perimeter calculation assuming 480p imgs
    height, width := 400, 800
    perimeter := 2 * (height + width)
    fmt.Println(perimeter, imgURL)

    return nil
}

// update status
func updateJobStatus(jobID string, status string, failedStoreID string) {
    jobsID, _ := primitive.ObjectIDFromHex(jobID)
    update := bson.M{"$set": bson.M{
        "status": status,
    }}
    if status == "failed" {
        update["$push"] = bson.M{
            "failed_id": failedStoreID,
        }
    }

    _, err := rabbitmq.DBCollection.UpdateOne(context.TODO(), bson.M{
        "_id": jobsID,
    }, update)

    if err != nil {
        log.Printf("Failed to update job status: %v", err)
    }
}