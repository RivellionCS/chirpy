package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// create a valid token to reuse
	validToken, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			tokenSecret: secret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong-secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "malformed token",
			tokenString: "this.is.not.valid",
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("got %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	tokenString, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected token string, got empty")
	}

	// validate it immediately
	gotUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("token failed validation: %v", err)
	}

	if gotUserID != userID {
		t.Errorf("got %v, want %v", gotUserID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// create a token that expires immediately
	tokenString, err := MakeJWT(userID, secret, 1*time.Nanosecond)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// ensure it's expired
	time.Sleep(10 * time.Millisecond)

	gotUserID, err := ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected error for expired token, got none")
	}

	if gotUserID != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %v", gotUserID)
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string
		headerValue string
		expected string
		expectError bool
	}{
		{
			name: "valid input",
			headerValue: "Bearer mytoken",
			expected: "mytoken",
			expectError: false,
		},
		{
			name: "empty input",
			headerValue: "",
			expected: "",
			expectError: true,
		},
		{
			name: "wrong input",
			headerValue: "Bearermytoken",
			expected: "",
			expectError: true,
		},
		{
			name: "empty token",
			headerValue: "Bearer ",
			expected: "",
			expectError: true,
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Authorization", tt.headerValue)
			result, err := GetBearerToken(headers)

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && result != tt.expected {
				t.Errorf("expected mytoken but got something else: %v", result)
			}
		})
	}
}