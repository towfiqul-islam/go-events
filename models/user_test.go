package models

import (
	"database/sql"
	"errors"
	"testing"

	"example.com/rest-api/test"
	"example.com/rest-api/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUser_Save(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testUser := test.GetTestUser()
	user := User{
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	tests := []struct {
		name    string
		user    User
		mockFn  func()
		wantErr bool
		wantID  int64
	}{
		{
			name: "Successful save",
			user: user,
			mockFn: func() {
				query := `INSERT INTO users\(email, password\) VALUES \(\?, \?\)`
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
			wantID:  1,
		},
		{
			name: "Prepare error",
			user: user,
			mockFn: func() {
				query := `INSERT INTO users\(email, password\) VALUES \(\?, \?\)`
				mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))
			},
			wantErr: true,
			wantID:  0,
		},
		{
			name: "Exec error",
			user: user,
			mockFn: func() {
				query := `INSERT INTO users\(email, password\) VALUES \(\?, \?\)`
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(errors.New("exec error"))
			},
			wantErr: true,
			wantID:  0,
		},
		{
			name: "LastInsertId error",
			user: user,
			mockFn: func() {
				query := `INSERT INTO users\(email, password\) VALUES \(\?, \?\)`
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

			err := tt.user.Save()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), tt.user.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, tt.user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUser_ValidateUser(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testUser := test.GetTestUser()
	plainPassword := "password123"
	hashedPassword, err := utils.HashPassword(plainPassword)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		user    User
		mockFn  func()
		wantErr bool
	}{
		{
			name: "Valid credentials",
			user: User{
				Email:    testUser.Email,
				Password: plainPassword,
			},
			mockFn: func() {
				columns := []string{"id", "password"}
				rows := sqlmock.NewRows(columns).AddRow(testUser.ID, hashedPassword)
				mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
					WithArgs(testUser.Email).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Invalid password",
			user: User{
				Email:    testUser.Email,
				Password: "wrongpassword",
			},
			mockFn: func() {
				columns := []string{"id", "password"}
				rows := sqlmock.NewRows(columns).AddRow(testUser.ID, hashedPassword)
				mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
					WithArgs(testUser.Email).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "User not found",
			user: User{
				Email:    "nonexistent@example.com",
				Password: plainPassword,
			},
			mockFn: func() {
				mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
					WithArgs("nonexistent@example.com").WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "Query error",
			user: User{
				Email:    testUser.Email,
				Password: plainPassword,
			},
			mockFn: func() {
				mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
					WithArgs(testUser.Email).WillReturnError(errors.New("query error"))
			},
			wantErr: true,
		},
		{
			name: "Scan error",
			user: User{
				Email:    testUser.Email,
				Password: plainPassword,
			},
			mockFn: func() {
				columns := []string{"id", "password"}
				rows := sqlmock.NewRows(columns).AddRow("invalid_id", hashedPassword)
				mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
					WithArgs(testUser.Email).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := tt.user.ValidateUser()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testUser.ID, tt.user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUser(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	testUser := test.GetTestUser()

	tests := []struct {
		name     string
		userID   int64
		mockFn   func()
		wantErr  bool
		wantUser *User
	}{
		{
			name:   "Successful query",
			userID: testUser.ID,
			mockFn: func() {
				columns := []string{"id", "email", "password"}
				rows := sqlmock.NewRows(columns).
					AddRow(testUser.ID, testUser.Email, testUser.Password)
				mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
					WithArgs(testUser.ID).WillReturnRows(rows)
			},
			wantErr: false,
			wantUser: &User{
				ID:       testUser.ID,
				Email:    testUser.Email,
				Password: testUser.Password,
			},
		},
		{
			name:   "User not found",
			userID: 999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
					WithArgs(int64(999)).WillReturnError(sql.ErrNoRows)
			},
			wantErr:  true,
			wantUser: nil,
		},
		{
			name:   "Query error",
			userID: testUser.ID,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
					WithArgs(testUser.ID).WillReturnError(errors.New("query error"))
			},
			wantErr:  true,
			wantUser: nil,
		},
		{
			name:   "Scan error",
			userID: testUser.ID,
			mockFn: func() {
				columns := []string{"id", "email", "password"}
				rows := sqlmock.NewRows(columns).
					AddRow("invalid_id", testUser.Email, testUser.Password)
				mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
					WithArgs(testUser.ID).WillReturnRows(rows)
			},
			wantErr:  true,
			wantUser: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			user, err := GetUser(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Email, user.Email)
				assert.Equal(t, tt.wantUser.Password, user.Password)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUser_ValidateUser_PasswordComparison(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	// Test with actual password hashing
	plainPassword := "testpassword123"
	hashedPassword, err := utils.HashPassword(plainPassword)
	assert.NoError(t, err)

	user := User{
		Email:    "test@example.com",
		Password: plainPassword,
	}

	// Mock successful database query
	columns := []string{"id", "password"}
	rows := sqlmock.NewRows(columns).AddRow(int64(1), hashedPassword)
	mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
		WithArgs(user.Email).WillReturnRows(rows)

	// Should successfully validate
	err = user.ValidateUser()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Save_Integration(t *testing.T) {
	mock, cleanup, err := test.SetupMockDB()
	assert.NoError(t, err)
	defer cleanup()

	user := User{
		Email:    "integration@example.com",
		Password: "hashedpassword",
	}

	// Mock successful save
	query := `INSERT INTO users\(email, password\) VALUES \(\?, \?\)`
	mock.ExpectPrepare(query).ExpectExec().
		WithArgs(user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(42, 1))

	err = user.Save()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
