package test

import (
	"database/sql/driver"
	"time"

	"example.com/rest-api/db"
	"github.com/DATA-DOG/go-sqlmock"
)

// SetupMockDB creates a mock database connection for testing
func SetupMockDB() (sqlmock.Sqlmock, func(), error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	// Replace the global DB with our mock
	originalDB := db.DB
	db.DB = mockDB

	// Return cleanup function
	cleanup := func() {
		mockDB.Close()
		db.DB = originalDB
	}

	return mock, cleanup, nil
}

// TestEvent represents a test event structure
type TestEvent struct {
	ID          int64
	Name        string
	Description string
	Location    string
	DateTime    time.Time
	UserID      int64
}

// GetTestEvent returns a sample event for testing
func GetTestEvent() TestEvent {
	return TestEvent{
		ID:          1,
		Name:        "Test Event",
		Description: "A test event description",
		Location:    "Test Location",
		DateTime:    time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC),
		UserID:      1,
	}
}

// TestUser represents a test user structure
type TestUser struct {
	ID       int64
	Email    string
	Password string
}

// GetTestUser returns a sample user for testing
func GetTestUser() TestUser {
	return TestUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword123",
	}
}

// TestNotification represents a test notification structure
type TestNotification struct {
	ID        int64
	UserID    int64
	EventID   int64
	Message   string
	Type      string
	IsRead    bool
	CreatedAt time.Time
}

// GetTestNotification returns a sample notification for testing
func GetTestNotification() TestNotification {
	return TestNotification{
		ID:        1,
		UserID:    1,
		EventID:   1,
		Message:   "Test notification message",
		Type:      "upcoming_event",
		IsRead:    false,
		CreatedAt: time.Now(),
	}
}

// TestEventRegister represents a test event registration structure
type TestEventRegister struct {
	ID      int64
	EventID int64
	UserID  int64
}

// GetTestEventRegister returns a sample event registration for testing
func GetTestEventRegister() TestEventRegister {
	return TestEventRegister{
		ID:      1,
		EventID: 1,
		UserID:  1,
	}
}

// ExpectExecSuccess sets up a mock expectation for a successful exec operation
func ExpectExecSuccess(mock sqlmock.Sqlmock, query string, lastInsertID int64, rowsAffected int64) {
	mock.ExpectPrepare(query).ExpectExec().WillReturnResult(
		sqlmock.NewResult(lastInsertID, rowsAffected))
}

// ExpectExecError sets up a mock expectation for a failed exec operation
func ExpectExecError(mock sqlmock.Sqlmock, query string, err error) {
	mock.ExpectPrepare(query).ExpectExec().WillReturnError(err)
}

// ExpectQueryRows sets up a mock expectation for a successful query operation
func ExpectQueryRows(mock sqlmock.Sqlmock, query string, columns []string, rows ...[]driver.Value) {
	mockRows := sqlmock.NewRows(columns)
	for _, row := range rows {
		mockRows.AddRow(row)
	}
	mock.ExpectQuery(query).WillReturnRows(mockRows)
}

// ExpectQueryError sets up a mock expectation for a failed query operation
func ExpectQueryError(mock sqlmock.Sqlmock, query string, err error) {
	mock.ExpectQuery(query).WillReturnError(err)
}

// ExpectPrepareError sets up a mock expectation for a failed prepare operation
func ExpectPrepareError(mock sqlmock.Sqlmock, query string, err error) {
	mock.ExpectPrepare(query).WillReturnError(err)
}
