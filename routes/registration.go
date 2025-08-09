package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func register(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id"})
		return
	}

	_, err = models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event"})
		fmt.Println(err)
		return
	}

	var EventRegister models.EventRegister

	EventRegister.EventID = eventId
	EventRegister.UserID = userId

	err = EventRegister.Register()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register in event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event registration success"})
}



func cancel(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id"})
		return
	}

	_, err = models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event"})
		fmt.Println(err)
		return
	}

	var EventRegister models.EventRegister

	EventRegister.EventID = eventId
	EventRegister.UserID = userId

	err = EventRegister.Cancel()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event cancelled"})
}