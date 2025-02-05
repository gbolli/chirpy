package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gbolli/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
        Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
    }

	type errorVal struct {
		Errormsg string `json:"error"`
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	w.Header().Set("Content-Type", "application/json")

	// Test chirp logic
	if len(params.Body) > 140 {
		respBody := errorVal {
			Errormsg: "Chirp is too long",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
    	w.Write(dat)
		return
	}

	// create chirp in DB
	newChirp := database.CreateChirpParams{
		Body: cleanBody(params.Body),
		UserID: params.UserID,
	}

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), newChirp)

	if err != nil {
		log.Printf("Error creating user in database: %s", err)
		w.WriteHeader(500)
		return
    }

	mainChirp := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	}

	// response

    dat, err := json.Marshal(mainChirp)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    
    w.WriteHeader(201)
    w.Write(dat)
}

func cleanBody(body string) string {
	bodyWords := strings.Fields(body)

	for i, word := range bodyWords {
		lowerword := strings.ToLower(word)
		if lowerword == "kerfuffle" || lowerword == "sharbert" || lowerword == "fornax" {
				bodyWords[i] = "****"
		}
	}

	return strings.Join(bodyWords, " ")
}