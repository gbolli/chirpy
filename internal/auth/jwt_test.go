package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTSuccess(t *testing.T) {
    userID := uuid.New()
    secret := "your-test-secret"
    // Create token that expires in 1 hour
    token, err := MakeJWT(userID, secret, time.Hour)
    if err != nil {
        t.Fatalf("Error creating token: %v", err)
    }
    
    // Validate the token
    gotID, err := ValidateJWT(token, secret)
    if err != nil {
        t.Fatalf("Error validating token: %v", err)
    }
    
    if gotID != userID {
        t.Errorf("got user ID %v, want %v", gotID, userID)
    }
}