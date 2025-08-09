package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	// events
	server.GET("/events", getEvents)
	server.GET("/events/:id", getSingleEvent)


	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", register)
	authenticated.DELETE("/events/:id/cancel", cancel)

	// users
	server.POST("/signup", signup)
	server.POST("/login", login)
	server.GET("/user/:id", getUserByID)
}