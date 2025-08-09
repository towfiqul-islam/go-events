package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/jobs"
	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func getNotifications(context *gin.Context) {
	userID, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	userId, ok := userID.(int64)
	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	notifications, err := models.GetNotificationsByUserID(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch notifications"})
		return
	}

	context.JSON(http.StatusOK, notifications)
}

func markNotificationAsRead(context *gin.Context) {
	notificationIDStr := context.Param("id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid notification ID"})
		return
	}

	userID, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	// Additional security: Verify the notification belongs to the user
	notifications, err := models.GetNotificationsByUserID(userID.(int64))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not verify notification ownership"})
		return
	}

	notificationExists := false
	for _, notification := range notifications {
		if notification.ID == notificationID {
			notificationExists = true
			break
		}
	}

	if !notificationExists {
		context.JSON(http.StatusForbidden, gin.H{"message": "Notification not found or access denied"})
		return
	}

	err = models.MarkNotificationAsRead(notificationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not mark notification as read"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func triggerNotificationCheck(context *gin.Context) {
	notificationService := jobs.NewNotificationService()
	err := notificationService.ProcessManually()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not process notifications"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Notification check triggered successfully"})
}
