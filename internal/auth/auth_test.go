package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if hash == "" {
		t.Fatalf("Expected non-empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Error checking password hash: %v", err)
	}
	if !match {
		t.Fatalf("Expected password to match hash")
	}

	wrongPassword := "wrongpassword"
	match, err = CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("Error checking password hash: %v", err)
	}
	if match {
		t.Fatalf("Expected wrong password not to match hash")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	match, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("Error verifying password: %v", err)
	}
	if !match {
		t.Fatalf("Expected password to match hash")
	}

	wrongPassword := "wrongpassword"
	match, err = VerifyPassword(wrongPassword, hash)
	if err != nil {
		t.Fatalf("Error verifying password: %v", err)
	}
	if match {
		t.Fatalf("Expected wrong password not to match hash")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	tokenSecret := "mysecret"
	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}
	token, err := MakeJWT(userID, tokenSecret, 1*time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}
	if validatedUserID != userID {
		t.Fatalf("Expected userID %s, got %s", userID, validatedUserID)
	}
}

func TestExpiredToken(t *testing.T) {
	tokenSecret := "mysecret"
	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}
	token, err := MakeJWT(userID, tokenSecret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("Expected error validating expired JWT, got nil")
	}
}

func TestWrongTokenSecret(t *testing.T) {
	tokenSecret := "mysecret"
	wrongTokenSecret := "wrongsecret"
	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}
	token, err := MakeJWT(userID, tokenSecret, 1*time.Hour)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongTokenSecret)
	if err == nil {
		t.Fatalf("Expected error validating JWT with wrong secret, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer mytoken"}
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Error getting bearer token: %v", err)
	}
	if token != "mytoken" {
		t.Fatalf("Expected token 'mytoken', got '%s'", token)
	}

	headers["Authorization"] = []string{"InvalidHeader"}
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Fatalf("Expected error for invalid header format, got nil")
	}

	delete(headers, "Authorization")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Fatalf("Expected error for missing authorization header, got nil")
	}
}
