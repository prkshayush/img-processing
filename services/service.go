package services

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prkshayush/img-processing/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// read-write mutex
var jobMap = struct {
    sync.RWMutex
    m map[string]*JobStatus
}{m: make(map[string]*JobStatus)}

type JobStatus struct {
    JobID       primitive.ObjectID
    Status      string
    FailedStore []string
    Mutex       sync.Mutex
    WaitGroup   sync.WaitGroup
}

// submit handler service
func HandleJobSubmit(request models.JobRequest) (string, error) {
    job := models.Job{
        Status:    "pending",
        CreatedAt: time.Now(),
        Visits:    request.Visits,
    }

    insertRes, err := models.InsertJob(job)
    if err != nil {
        return "", err
    }

    jobID := insertRes.InsertedID.(primitive.ObjectID)
    jobStatus := &JobStatus{
        JobID:  jobID,
        Status: "pending",
    }

    jobMap.Lock()
    jobMap.m[jobID.Hex()] = jobStatus
    jobMap.Unlock()

    go processJob(jobID.Hex(), request.Visits)

    return jobID.Hex(), nil
}

func processJob(jobID string, visits []models.Visit) {
    jobMap.RLock()
    jobStatus := jobMap.m[jobID]
    jobMap.RUnlock()

    for _, visit := range visits {
        for _, imageURL := range visit.ImageURLs {
            jobStatus.WaitGroup.Add(1)
            go processImage(jobID, visit.StoreID, imageURL)
        }
    }

    jobStatus.WaitGroup.Wait()

    jobMap.Lock()
    if len(jobStatus.FailedStore) > 0 {
        jobStatus.Status = "failed"
    } else {
        jobStatus.Status = "complete"
    }
    models.UpdateJobStatus(jobStatus.JobID, jobStatus.Status, jobStatus.FailedStore)
    delete(jobMap.m, jobID)
    jobMap.Unlock()
}

func processImage(jobID, storeID, imageURL string) {
    defer jobMap.m[jobID].WaitGroup.Done()

    response, err := http.Get(imageURL)
    if err != nil {
        updateFailedStore(jobID, storeID)
        return
    }
    defer response.Body.Close()

    tmpFile, err := os.CreateTemp("", "image-*.jpg")
    if err != nil {
        updateFailedStore(jobID, storeID)
        return
    }
    defer os.Remove(tmpFile.Name())

    _, err = io.Copy(tmpFile, response.Body)
    if err != nil {
        updateFailedStore(jobID, storeID)
        return
    }

    time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)

    // perimeter calculation assuming 480p images
    height, width := 400, 800
    perimeter := 2 * (height + width)

	// save results
    updateImageResult(jobID, storeID, imageURL, perimeter)
}

func updateImageResult(jobID, storeID, imageURL string, perimeter int) {
    jobObjectID, _ := primitive.ObjectIDFromHex(jobID)
    filter := bson.M{"_id": jobObjectID, "visits.store_id": storeID}
    update := bson.M{
        "$push": bson.M{
            "visits.$.results": models.ImageResult{
                ImageURL:  imageURL,
                Perimeter: perimeter,
                Processed: true,
            },
        },
    }
    _, err := models.DBCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Printf("Failed to update image result: %v", err)
    }
}

func updateFailedStore(jobID, storeID string) {
    jobMap.Lock()
    jobMap.m[jobID].FailedStore = append(jobMap.m[jobID].FailedStore, storeID)
    jobMap.Unlock()
}