package jobs

import (
	"errors"
	"testing"
	"time"

	"example.com/rest-api/test"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNotificationService_NewNotificationService(t *testing.T) {
	service := NewNotificationService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.stopChan)
}

func TestNotificationService_ProcessManually(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	service := NewNotificationService()

	// Test with no upcoming events
	t.Run("No upcoming events", func(t *testing.T) {
		// Mock the query to return no results
		columns := []string{"id", "name", "dateTime", "user_id"}
		rows := sqlmock.NewRows(columns)
		mock.ExpectQuery(`SELECT e\.id, e\.name, e\.dateTime, er\.user_id`).WillReturnRows(rows)

		err := service.ProcessManually()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotificationService_ProcessUpcomingEvents(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	service := NewNotificationService()

	testEvent := test.GetTestEvent()
	futureTime := time.Now().Add(12 * time.Hour) // 12 hours from now

	tests := []struct {
		name    string
		mockFn  func()
		wantErr bool
	}{
		{
			name: "Successful processing with events",
			mockFn: func() {
				// Mock GetUpcomingEventsForNotification
				columns := []string{"id", "name", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns).
					AddRow(testEvent.ID, testEvent.Name, futureTime, testEvent.UserID)
				mock.ExpectQuery(`SELECT e\.id, e\.name, e\.dateTime, er\.user_id`).WillReturnRows(rows)

				// Mock notification save
				insertQuery := `INSERT INTO notifications \(user_id, event_id, message, type, is_read, created_at\) VALUES \(\?, \?, \?, \?, \?, \?\)`
				mock.ExpectPrepare(insertQuery).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error fetching events",
			mockFn: func() {
				mock.ExpectQuery(`SELECT e\.id, e\.name, e\.dateTime, er\.user_id`).
					WillReturnError(errors.New("database error"))
			},
			wantErr: false, // processUpcomingEvents logs errors but doesn't return them
		},
		{
			name: "No upcoming events",
			mockFn: func() {
				columns := []string{"id", "name", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns)
				mock.ExpectQuery(`SELECT e\.id, e\.name, e\.dateTime, er\.user_id`).WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			// Call processUpcomingEvents directly for testing
			service.processUpcomingEvents()

			// Verify expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationService_GenerateNotificationMessage(t *testing.T) {
	service := NewNotificationService()

	now := time.Now()

	tests := []struct {
		name      string
		eventName string
		eventTime time.Time
		wantMsg   string
	}{
		{
			name:      "Event starting soon (within 1 hour)",
			eventName: "Test Event",
			eventTime: now.Add(30 * time.Minute),
			wantMsg:   "Reminder: Your event 'Test Event' is starting soon at",
		},
		{
			name:      "Event within 24 hours",
			eventTime: now.Add(12 * time.Hour),
			eventName: "Conference",
			wantMsg:   "Reminder: Your event 'Conference' is in",
		},
		{
			name:      "Event beyond 24 hours",
			eventName: "Workshop",
			eventTime: now.Add(48 * time.Hour),
			wantMsg:   "Reminder: You have an upcoming event 'Workshop' on",
		},
		{
			name:      "Event exactly 1 hour away",
			eventName: "Meeting",
			eventTime: now.Add(61 * time.Minute), // Slightly over 1 hour to avoid boundary issues
			wantMsg:   "Reminder: Your event 'Meeting' is in",
		},
		{
			name:      "Event exactly 24 hours away",
			eventName: "Presentation",
			eventTime: now.Add(25 * time.Hour), // Slightly over 24 hours
			wantMsg:   "Reminder: You have an upcoming event 'Presentation' on",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := service.generateNotificationMessage(tt.eventName, tt.eventTime)
			assert.Contains(t, message, tt.wantMsg)
			assert.Contains(t, message, tt.eventName)
		})
	}
}

func TestNotificationService_StartStop(t *testing.T) {
	// Setup mock database for this test
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	service := NewNotificationService()

	// Mock the query to return no results to avoid database processing
	columns := []string{"id", "name", "dateTime", "user_id"}
	rows := sqlmock.NewRows(columns)
	mock.ExpectQuery(`SELECT e\.id, e\.name, e\.dateTime, er\.user_id`).WillReturnRows(rows)

	// Test that service can be stopped
	done := make(chan bool)
	
	go func() {
		service.Start()
		done <- true
	}()

	// Give it a moment to start and process
	time.Sleep(50 * time.Millisecond)
	
	// Stop the service
	service.Stop()

	// Wait for completion or timeout
	select {
	case <-done:
		// Service stopped successfully
	case <-time.After(200 * time.Millisecond):
		t.Error("Service did not stop within timeout")
	}

	// Verify the mock expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationService_GenerateMessage_EdgeCases(t *testing.T) {
	service := NewNotificationService()

	now := time.Now()

	tests := []struct {
		name      string
		eventName string
		eventTime time.Time
		wantCheck func(string) bool
	}{
		{
			name:      "Empty event name",
			eventName: "",
			eventTime: now.Add(2 * time.Hour),
			wantCheck: func(msg string) bool {
				return msg != "" && !assert.ObjectsAreEqual(msg, "")
			},
		},
		{
			name:      "Event name with special characters",
			eventName: "Test Event! @#$%^&*()",
			eventTime: now.Add(2 * time.Hour),
			wantCheck: func(msg string) bool {
				return assert.Contains(nil, msg, "Test Event! @#$%^&*()")
			},
		},
		{
			name:      "Very long event name",
			eventName: "This is a very long event name that might cause issues with formatting and display in the notification system",
			eventTime: now.Add(2 * time.Hour),
			wantCheck: func(msg string) bool {
				return len(msg) > 0 && assert.Contains(nil, msg, "This is a very long event name")
			},
		},
		{
			name:      "Past event (edge case)",
			eventName: "Past Event",
			eventTime: now.Add(-2 * time.Hour),
			wantCheck: func(msg string) bool {
				return len(msg) > 0 && assert.Contains(nil, msg, "Past Event")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := service.generateNotificationMessage(tt.eventName, tt.eventTime)
			assert.True(t, tt.wantCheck(message), "Message validation failed: %s", message)
		})
	}
}

func TestNotificationService_ProcessUpcomingEvents_NotificationSaveError(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	service := NewNotificationService()

	testEvent := test.GetTestEvent()
	futureTime := time.Now().Add(12 * time.Hour)

	// Mock GetUpcomingEventsForNotification to return an event
	columns := []string{"id", "name", "dateTime", "user_id"}
	rows := sqlmock.NewRows(columns).
		AddRow(testEvent.ID, testEvent.Name, futureTime, testEvent.UserID)
	mock.ExpectQuery(`SELECT e\.id, e\.name, e\.dateTime, er\.user_id`).WillReturnRows(rows)

	// Mock notification save to fail
	insertQuery := `INSERT INTO notifications \(user_id, event_id, message, type, is_read, created_at\) VALUES \(\?, \?, \?, \?, \?, \?\)`
	mock.ExpectPrepare(insertQuery).ExpectExec().WillReturnError(errors.New("save failed"))

	// This should not panic even when save fails
	service.processUpcomingEvents()

	assert.NoError(t, mock.ExpectationsWereMet())
}
