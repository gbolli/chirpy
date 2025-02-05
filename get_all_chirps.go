package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	
	// get all chirps

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	allChirps, err := cfg.dbQueries.GetAllChirps(r.Context())

	if err != nil {
		log.Printf("Error getting all chirps in database: %s", err)
		w.WriteHeader(500)
		return
    }

	var mainChirps []Chirp;
	for _, chirp := range allChirps {
		mainChirps = append(
			mainChirps,
			Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
			},
		) 
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