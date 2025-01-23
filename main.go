package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prkshayush/img-processing/rabbitmq"
	"github.com/prkshayush/img-processing/routes"
	workers "github.com/prkshayush/img-processing/utils"
)

func main() {
	fmt.Println("Yo, project")

	err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to establish connection with RabbitMQ: %v", err)
	}

	// workers as go-routine
	go workers.ProcessJobs()

	// router
	r := gin.Default()
	routes.ApiRoutes(r)

	log.Println("Starting server on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}