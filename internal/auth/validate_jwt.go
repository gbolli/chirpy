package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    	// Check if the signing method is what we expect (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		// Return the secret key as a []byte
		return []byte(tokenSecret), nil
	})
	if err != nil { return uuid.Nil, err }
	// Check if token is valid
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	// Get the subject claim
	subjectID := claims.Subject

	// Convert string to UUID
	userID, err := uuid.Parse(subjectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}