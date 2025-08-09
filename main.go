package main

import (
	"example.com/rest-api/db"
	"example.com/rest-api/jobs"
	"example.com/rest-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	// Start the notification service
	notificationService := jobs.NewNotificationService()
	notificationService.Start()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")

}
