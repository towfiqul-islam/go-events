package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Generate a valid token for testing
	email := "test@example.com"
	userID := int64(123)
	validToken, err := utils.GenerateToken(email, userID)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldAbort    bool
		expectedUserID int64
	}{
		{
			name:           "Valid token",
			authHeader:     validToken,
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
			expectedUserID: userID,
		},
		{
			name:           "Empty Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
			expectedUserID: 0,
		},
		{
			name:           "Invalid token format",
			authHeader:     "invalid.token.format",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
			expectedUserID: 0,
		},
		{
			name:           "Malformed token",
			authHeader:     "notavalidtoken",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
			expectedUserID: 0,
		},
		{
			name:           "Token with wrong signature",
			authHeader:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyMywiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjcwMDAwMDAwfQ.wrongsignature",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
			expectedUserID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin router for each test
			router := gin.New()

			// Add the middleware
			router.Use(Authenticate)

			// Add a test route that should only be accessible with valid auth
			router.GET("/test", func(c *gin.Context) {
				userID, exists := c.Get("userId")
				if exists {
					c.JSON(http.StatusOK, gin.H{"userId": userID})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "no userId"})
				}
			})

			// Create a test request
			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			// Set the Authorization header if provided
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldAbort {
				// Should return unauthorized response
				assert.Contains(t, w.Body.String(), "Not authorized")
			} else {
				// Should pass through to the handler and have userId set
				assert.Contains(t, w.Body.String(), "userId")
			}
		})
	}
}

func TestAuthenticate_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a complete flow test
	email := "integration@example.com"
	userID := int64(456)
	token, err := utils.GenerateToken(email, userID)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(Authenticate)

	// Handler that checks if userId is properly set
	router.GET("/protected", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		assert.True(t, exists, "userId should be set in context")
		assert.Equal(t, userID, userId.(int64), "userId should match the token")
		c.JSON(http.StatusOK, gin.H{"success": true, "userId": userId})
	})

	req, err := http.NewRequest("GET", "/protected", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "456")
}

func TestAuthenticate_ContextValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	email := "context@example.com"
	userID := int64(789)
	token, err := utils.GenerateToken(email, userID)
	assert.NoError(t, err)

	var capturedUserID int64
	var contextExists bool

	router := gin.New()
	router.Use(Authenticate)

	router.GET("/capture", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		contextExists = exists
		if exists {
			capturedUserID = userId.(int64)
		}
		c.Status(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/capture", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, contextExists, "Context should contain userId")
	assert.Equal(t, userID, capturedUserID, "Captured userID should match")
}

func TestAuthenticate_Next(t *testing.T) {
	gin.SetMode(gin.TestMode)

	email := "next@example.com"
	userID := int64(999)
	token, err := utils.GenerateToken(email, userID)
	assert.NoError(t, err)

	nextCalled := false

	router := gin.New()
	router.Use(Authenticate)
	router.Use(func(c *gin.Context) {
		nextCalled = true
		c.Next()
	})

	router.GET("/next", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/next", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, nextCalled, "Next middleware should be called when authentication succeeds")
}
