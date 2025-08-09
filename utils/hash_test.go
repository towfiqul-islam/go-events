package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "Long password",
			password: "verylongpasswordthatexceedsnormallengthbutshouldsitllwork123456789",
			wantErr:  false,
		},
		{
			name:     "Special characters password",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.NotEqual(t, tt.password, hashedPassword)

				// Verify the hash starts with bcrypt prefix
				assert.Contains(t, hashedPassword, "$2a$")
			}
		})
	}
}

func TestCheckHashPassword(t *testing.T) {
	// First, create a known hash
	plainPassword := "testpassword123"
	hashedPassword, err := HashPassword(plainPassword)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		want           bool
	}{
		{
			name:           "Correct password",
			password:       plainPassword,
			hashedPassword: hashedPassword,
			want:           true,
		},
		{
			name:           "Incorrect password",
			password:       "wrongpassword",
			hashedPassword: hashedPassword,
			want:           false,
		},
		{
			name:           "Empty password with valid hash",
			password:       "",
			hashedPassword: hashedPassword,
			want:           false,
		},
		{
			name:           "Valid password with empty hash",
			password:       plainPassword,
			hashedPassword: "",
			want:           false,
		},
		{
			name:           "Both empty",
			password:       "",
			hashedPassword: "",
			want:           false,
		},
		{
			name:           "Invalid hash format",
			password:       plainPassword,
			hashedPassword: "invalidhash",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckHashPassword(tt.password, tt.hashedPassword)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestHashPassword_Consistency(t *testing.T) {
	password := "testpassword"

	// Hash the same password multiple times
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2) // Should be different due to salt

	// But both should validate the original password
	assert.True(t, CheckHashPassword(password, hash1))
	assert.True(t, CheckHashPassword(password, hash2))
}

func TestHashPassword_EmptyString(t *testing.T) {
	hashedEmpty, err := HashPassword("")
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedEmpty)

	// Empty password should validate against its hash
	assert.True(t, CheckHashPassword("", hashedEmpty))

	// But not against any other password
	assert.False(t, CheckHashPassword("notempty", hashedEmpty))
}
