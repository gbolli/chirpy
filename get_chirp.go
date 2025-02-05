package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	
	// chirp ID

	chirpID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(chirpID)

	if err != nil {
    // Handle invalid UUID, likely return a 404 or 400
    	http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), parsedID)

	if err != nil {
		log.Printf("Error getting chirp from database: %s", err)
		w.WriteHeader(404)
		return
    }

	mainChirps :=
			Chirp{
				ID: dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				Body: dbChirp.Body,
				UserID: dbChirp.UserID,
			}

	// response

    dat, err := json.Marshal(mainChirps)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    
    w.WriteHeader(200)
    w.Write(dat)

}