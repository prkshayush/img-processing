package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prkshayush/img-processing/models"
	"github.com/prkshayush/img-processing/routes"
	"github.com/prkshayush/img-processing/utils"
)

func main() {
    // environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    mongoURI := os.Getenv("MONGODB_URI")
    dbName := os.Getenv("MONGODB_DB")
    collectionName := os.Getenv("MONGODB_COLLECTION")
    masterStorePath := os.Getenv("MASTER_STORE_PATH")

    err = models.ConnectDB(mongoURI, dbName, collectionName)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    err = utils.LoadMasterStore(masterStorePath)
    if err != nil {
        log.Fatalf("Failed to laod master data: %v", err)
    }

    router := gin.Default()
    routes.ApiRoutes(router)

    log.Println("Starting server on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}