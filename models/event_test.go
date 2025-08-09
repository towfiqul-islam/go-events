package models

import (
	"database/sql"
	"errors"
	"testing"

	"example.com/rest-api/test"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestEvent_Save(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testEvent := test.GetTestEvent()
	event := Event{
		Name:        testEvent.Name,
		Description: testEvent.Description,
		Location:    testEvent.Location,
		DateTime:    testEvent.DateTime,
		UserID:      testEvent.UserID,
	}

	tests := []struct {
		name    string
		event   Event
		mockFn  func()
		wantErr bool
		wantID  int64
	}{
		{
			name:  "Successful save",
			event: event,
			mockFn: func() {
				query := `INSERT INTO events \(name, description, location, dateTime, user_id\) VALUES \(\?, \?, \?, \?, \?\)`
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
			wantID:  1,
		},
		{
			name:  "Prepare error",
			event: event,
			mockFn: func() {
				query := `INSERT INTO events \(name, description, location, dateTime, user_id\) VALUES \(\?, \?, \?, \?, \?\)`
				mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))
			},
			wantErr: true,
			wantID:  0,
		},
		{
			name:  "Exec error",
			event: event,
			mockFn: func() {
				query := `INSERT INTO events \(name, description, location, dateTime, user_id\) VALUES \(\?, \?, \?, \?, \?\)`
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(errors.New("exec error"))
			},
			wantErr: true,
			wantID:  0,
		},
		{
			name:  "LastInsertId error",
			event: event,
			mockFn: func() {
				query := `INSERT INTO events \(name, description, location, dateTime, user_id\) VALUES \(\?, \?, \?, \?, \?\)`
				result := sqlmock.NewErrorResult(errors.New("last insert id error"))
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(result)
			},
			wantErr: true,
			wantID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := tt.event.Save()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), tt.event.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, tt.event.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEvent_Update(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testEvent := test.GetTestEvent()
	event := Event{
		ID:          testEvent.ID,
		Name:        testEvent.Name,
		Description: testEvent.Description,
		Location:    testEvent.Location,
		DateTime:    testEvent.DateTime,
		UserID:      testEvent.UserID,
	}

	tests := []struct {
		name    string
		event   Event
		mockFn  func()
		wantErr bool
	}{
		{
			name:  "Successful update",
			event: event,
			mockFn: func() {
				query := `UPDATE events SET name = \?, description = \?, location = \?, dateTime = \? WHERE id = \?`
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:  "Prepare error",
			event: event,
			mockFn: func() {
				query := `UPDATE events SET name = \?, description = \?, location = \?, dateTime = \? WHERE id = \?`
				mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))
			},
			wantErr: true,
		},
		{
			name:  "Exec error",
			event: event,
			mockFn: func() {
				query := `UPDATE events SET name = \?, description = \?, location = \?, dateTime = \? WHERE id = \?`
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(errors.New("exec error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := tt.event.Update()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEvent_Delete(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testEvent := test.GetTestEvent()
	event := Event{
		ID: testEvent.ID,
	}

	tests := []struct {
		name    string
		event   Event
		mockFn  func()
		wantErr bool
	}{
		{
			name:  "Successful delete",
			event: event,
			mockFn: func() {
				query := `DELETE FROM events WHERE id = \?`
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:  "Prepare error",
			event: event,
			mockFn: func() {
				query := `DELETE FROM events WHERE id = \?`
				mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))
			},
			wantErr: true,
		},
		{
			name:  "Exec error",
			event: event,
			mockFn: func() {
				query := `DELETE FROM events WHERE id = \?`
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(errors.New("exec error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := tt.event.Delete()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllEvents(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testEvent := test.GetTestEvent()

	tests := []struct {
		name       string
		mockFn     func()
		wantErr    bool
		wantCount  int
		wantEvents []Event
	}{
		{
			name: "Successful query with results",
			mockFn: func() {
				columns := []string{"id", "name", "description", "location", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns).
					AddRow(testEvent.ID, testEvent.Name, testEvent.Description,
						testEvent.Location, testEvent.DateTime, testEvent.UserID)
				mock.ExpectQuery(`SELECT \* FROM events`).WillReturnRows(rows)
			},
			wantErr:   false,
			wantCount: 1,
			wantEvents: []Event{
				{
					ID:          testEvent.ID,
					Name:        testEvent.Name,
					Description: testEvent.Description,
					Location:    testEvent.Location,
					DateTime:    testEvent.DateTime,
					UserID:      testEvent.UserID,
				},
			},
		},
		{
			name: "Successful query with no results",
			mockFn: func() {
				columns := []string{"id", "name", "description", "location", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns)
				mock.ExpectQuery(`SELECT \* FROM events`).WillReturnRows(rows)
			},
			wantErr:    false,
			wantCount:  0,
			wantEvents: []Event{},
		},
		{
			name: "Query error",
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM events`).WillReturnError(errors.New("query error"))
			},
			wantErr:    true,
			wantCount:  0,
			wantEvents: nil,
		},
		{
			name: "Scan error",
			mockFn: func() {
				columns := []string{"id", "name", "description", "location", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns).
					AddRow("invalid_id", testEvent.Name, testEvent.Description,
						testEvent.Location, testEvent.DateTime, testEvent.UserID)
				mock.ExpectQuery(`SELECT \* FROM events`).WillReturnRows(rows)
			},
			wantErr:    true,
			wantCount:  0,
			wantEvents: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			events, err := GetAllEvents()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, events)
			} else {
				assert.NoError(t, err)
				assert.Len(t, events, tt.wantCount)
				if tt.wantCount > 0 {
					assert.Equal(t, tt.wantEvents[0].ID, events[0].ID)
					assert.Equal(t, tt.wantEvents[0].Name, events[0].Name)
					assert.Equal(t, tt.wantEvents[0].Description, events[0].Description)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetEventById(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testEvent := test.GetTestEvent()

	tests := []struct {
		name      string
		eventID   int64
		mockFn    func()
		wantErr   bool
		wantEvent *Event
	}{
		{
			name:    "Successful query",
			eventID: testEvent.ID,
			mockFn: func() {
				columns := []string{"id", "name", "description", "location", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns).
					AddRow(testEvent.ID, testEvent.Name, testEvent.Description,
						testEvent.Location, testEvent.DateTime, testEvent.UserID)
				mock.ExpectQuery(`SELECT \* FROM events WHERE id = \?`).
					WithArgs(testEvent.ID).WillReturnRows(rows)
			},
			wantErr: false,
			wantEvent: &Event{
				ID:          testEvent.ID,
				Name:        testEvent.Name,
				Description: testEvent.Description,
				Location:    testEvent.Location,
				DateTime:    testEvent.DateTime,
				UserID:      testEvent.UserID,
			},
		},
		{
			name:    "Event not found",
			eventID: 999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM events WHERE id = \?`).
					WithArgs(int64(999)).WillReturnError(sql.ErrNoRows)
			},
			wantErr:   true,
			wantEvent: nil,
		},
		{
			name:    "Query error",
			eventID: testEvent.ID,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM events WHERE id = \?`).
					WithArgs(testEvent.ID).WillReturnError(errors.New("query error"))
			},
			wantErr:   true,
			wantEvent: nil,
		},
		{
			name:    "Scan error",
			eventID: testEvent.ID,
			mockFn: func() {
				columns := []string{"id", "name", "description", "location", "dateTime", "user_id"}
				rows := sqlmock.NewRows(columns).
					AddRow("invalid_id", testEvent.Name, testEvent.Description,
						testEvent.Location, testEvent.DateTime, testEvent.UserID)
				mock.ExpectQuery(`SELECT \* FROM events WHERE id = \?`).
					WithArgs(testEvent.ID).WillReturnRows(rows)
			},
			wantErr:   true,
			wantEvent: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			event, err := GetEventById(tt.eventID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.wantEvent.ID, event.ID)
				assert.Equal(t, tt.wantEvent.Name, event.Name)
				assert.Equal(t, tt.wantEvent.Description, event.Description)
				assert.Equal(t, tt.wantEvent.Location, event.Location)
				assert.Equal(t, tt.wantEvent.UserID, event.UserID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
