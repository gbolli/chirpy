package main

import (
	"log"
	"net/http"

	"github.com/gbolli/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	
	// chirp ID

	chirpID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(chirpID)

	if err != nil {
    // Handle invalid UUID, likely return a 404 or 400
    	http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	// Authenticate user token

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(401)
		return
	}

	tokenUserID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(401)
		return
	}

	// Check for correct owner of chirp

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), parsedID)
	if err != nil {
		log.Printf("Error getting chirp from database: %s", err)
		w.WriteHeader(404)
		return
	}
	if dbChirp.UserID != tokenUserID {
		log.Printf("Unauthorized user: %s", err)
		w.WriteHeader(403)
		return
	}

	// Delete chirp in database

	err = cfg.dbQueries.DeleteChirp(r.Context(), parsedID)
	if err != nil {
		log.Printf("Error getting chirp from database: %s", err)
		w.WriteHeader(404)
		return
    }

    w.WriteHeader(204)
}