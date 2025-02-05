package auth

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
    password := "mySecretPassword123"
    
    // Test that hashing works
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    if hash == password {
        t.Error("Hash should not be equal to original password")
    }
    
    // Test that same password hashes differently (due to salt)
    hash2, err := HashPassword(password)
    if err != nil {
        t.Fatalf("Failed to hash password second time: %v", err)
    }
    if hash == hash2 {
        t.Error("Hashes should be different due to random salt")
    }
}

func TestCheckPasswordHash(t *testing.T) {
    password := "mySecretPassword123"
    
    // First create a hash
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    
    // Test correct password
    err = CheckPasswordHash(password, hash)
    if err != nil {
		fmt.Print(err)
        t.Error("Password should match hash")
    }
    
    // Test incorrect password
    wrongPassword := "wrongPassword123"
    err = CheckPasswordHash(wrongPassword, hash)
    if err == nil {
        t.Error("Wrong password should not match hash")
    }
}