package jobs

import (
	"fmt"
	"log"
	"time"

	"example.com/rest-api/models"
)

// NotificationService handles background job for notifications
type NotificationService struct {
	stopChan chan bool
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		stopChan: make(chan bool),
	}
}

// Start begins the background job that runs every hour
func (ns *NotificationService) Start() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		log.Println("Notification service started")

		// Run immediately when started
		ns.processUpcomingEvents()

		for {
			select {
			case <-ticker.C:
				ns.processUpcomingEvents()
			case <-ns.stopChan:
				ticker.Stop()
				log.Println("Notification service stopped")
				return
			}
		}
	}()
}

// Stop stops the background job
func (ns *NotificationService) Stop() {
	ns.stopChan <- true
}

// processUpcomingEvents finds upcoming events and creates notifications
func (ns *NotificationService) processUpcomingEvents() {
	log.Println("Processing upcoming events for notifications...")

	upcomingEvents, err := models.GetUpcomingEventsForNotification()
	if err != nil {
		log.Printf("Error fetching upcoming events: %v", err)
		return
	}

	if len(upcomingEvents) == 0 {
		log.Println("No upcoming events found for notifications")
		return
	}

	notificationsCreated := 0

	for _, event := range upcomingEvents {
		notification := &models.Notification{
			UserID:    event.UserID,
			EventID:   event.EventID,
			Message:   ns.generateNotificationMessage(event.EventName, event.DateTime),
			Type:      "upcoming_event",
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		err := notification.Save()
		if err != nil {
			log.Printf("Error creating notification for user %d, event %d: %v",
				event.UserID, event.EventID, err)
			continue
		}

		notificationsCreated++
		log.Printf("Created notification for user %d for event '%s'",
			event.UserID, event.EventName)
	}

	log.Printf("Successfully created %d notifications for upcoming events", notificationsCreated)
}

// generateNotificationMessage creates a user-friendly notification message
func (ns *NotificationService) generateNotificationMessage(eventName string, eventTime time.Time) string {
	hoursUntil := time.Until(eventTime).Hours()

	if hoursUntil <= 1 {
		return fmt.Sprintf("Reminder: Your event '%s' is starting soon at %s!",
			eventName, eventTime.Format("3:04 PM"))
	} else if hoursUntil <= 24 {
		hours := int(hoursUntil)
		return fmt.Sprintf("Reminder: Your event '%s' is in %d hour(s) at %s",
			eventName, hours, eventTime.Format("3:04 PM on Jan 2"))
	}

	return fmt.Sprintf("Reminder: You have an upcoming event '%s' on %s",
		eventName, eventTime.Format("January 2, 2006 at 3:04 PM"))
}

// ProcessManually allows manual triggering of the notification process
// This can be useful for testing or manual runs
func (ns *NotificationService) ProcessManually() error {
	log.Println("Manual notification processing triggered")
	ns.processUpcomingEvents()
	return nil
}
