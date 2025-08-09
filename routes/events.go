package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getSingleEvent(context *gin.Context) {
	eventId := context.Param("id")

	id, err := strconv.ParseInt(eventId, 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse eventID"})
		return
	}

	event, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find event"})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {

	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		fmt.Println(err)
		return
	}

	event.UserID = context.GetInt64("userId")

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "event created", "event": event})

}

func updateEvent(context *gin.Context) {

	eventId := context.Param("id")

	id, err := strconv.ParseInt(eventId, 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse eventID"})
		return
	}

	event, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find event"})
		return
	}

	userId := context.GetInt64("userId")

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized to update event"})
		return
	}

	err = context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		fmt.Println(err)
		return
	}

	event.ID = id

	err = event.Update()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event", "error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event updated"})

}

func deleteEvent(context *gin.Context) {

	eventId := context.Param("id")

	id, err := strconv.ParseInt(eventId, 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse eventID"})
		return
	}

	event, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find event"})
		return
	}

	err = event.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete event"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event deleted"})

}
