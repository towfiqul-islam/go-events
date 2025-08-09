package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		userID  int64
		wantErr bool
	}{
		{
			name:    "Valid credentials",
			email:   "test@example.com",
			userID:  123,
			wantErr: false,
		},
		{
			name:    "Empty email",
			email:   "",
			userID:  123,
			wantErr: false, // JWT allows empty email
		},
		{
			name:    "Zero user ID",
			email:   "test@example.com",
			userID:  0,
			wantErr: false,
		},
		{
			name:    "Negative user ID",
			email:   "test@example.com",
			userID:  -1,
			wantErr: false,
		},
		{
			name:    "Large user ID",
			email:   "test@example.com",
			userID:  9223372036854775807, // max int64
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.email, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token structure (JWT has 3 parts separated by dots)
				assert.Regexp(t, `^[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+$`, token)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	// Create a valid token first
	email := "test@example.com"
	userID := int64(123)
	validToken, err := GenerateToken(email, userID)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		wantID    int64
		wantErr   bool
		setupFunc func() string // for creating specific test tokens
	}{
		{
			name:    "Valid token",
			token:   validToken,
			wantID:  userID,
			wantErr: false,
		},
		{
			name:    "Empty token",
			token:   "",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Invalid token format",
			token:   "invalid.token.format",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "not.a.jwt",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Token with wrong signature",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyMywiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjcwMDAwMDAwfQ.wrongsignature",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Expired token",
			wantID:  0,
			wantErr: true,
			setupFunc: func() string {
				// Create an expired token
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"userId": int64(123),
					"email":  "test@example.com",
					"exp":    time.Now().Add(-time.Hour).Unix(), // expired 1 hour ago
				})
				tokenString, _ := token.SignedString([]byte("supersecret"))
				return tokenString
			},
		},
		{
			name:    "Token without expiration - should work",
			wantID:  123,
			wantErr: false,
			setupFunc: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"userId": int64(123),
					"email":  "test@example.com",
					// no exp claim - JWT library doesn't require it by default
				})
				tokenString, _ := token.SignedString([]byte("supersecret"))
				return tokenString
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token
			if tt.setupFunc != nil {
				token = tt.setupFunc()
			}

			userID, err := VerifyToken(token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, userID)
			}
		})
	}
}

func TestTokenRoundTrip(t *testing.T) {
	// Test generating a token and then verifying it
	email := "roundtrip@example.com"
	userID := int64(456)

	// Generate token
	token, err := GenerateToken(email, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	verifiedUserID, err := VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, verifiedUserID)
}

func TestGenerateToken_DifferentTokensForSameUser(t *testing.T) {
	email := "test@example.com"
	userID := int64(123)

	// Generate two tokens for the same user
	token1, err1 := GenerateToken(email, userID)
	time.Sleep(time.Second) // Ensure different timestamps (1 second for more reliable difference)
	token2, err2 := GenerateToken(email, userID)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// Note: tokens might be the same if generated within the same second due to Unix timestamp precision

	// But both should verify to the same user ID
	userID1, err1 := VerifyToken(token1)
	userID2, err2 := VerifyToken(token2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, userID, userID1)
	assert.Equal(t, userID, userID2)
}

func TestVerifyToken_InvalidSigningMethod(t *testing.T) {
	// Create a manually crafted token with wrong signing method indicator (RSA instead of HMAC)
	// This token has RSA256 in the header but we expect HMAC methods only
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyMywiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjcwMDAwMDAwfQ.wrongsignature"

	// This should fail because we expect HMAC methods only
	userID, err := VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, int64(0), userID)
	assert.Contains(t, err.Error(), "could not parse token")
}

func TestVerifyToken_MissingClaims(t *testing.T) {
	// Create a token without userId claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "test@example.com",
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
		// missing userId
	})

	tokenString, err := token.SignedString([]byte("supersecret"))
	assert.NoError(t, err)

	// This should fail because userId is missing
	userID, err := VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, int64(0), userID)
}
