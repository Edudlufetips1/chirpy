package auth

import (
	"testing"
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
