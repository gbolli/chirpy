package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)


func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type email struct {
        Email string `json:"email"`
    }

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
    em := email{}
    err := decoder.Decode(&em)

    if err != nil {
		log.Printf("Error decoding email: %s", err)
		w.WriteHeader(500)
		return
    }

	nullEmail := sql.NullString{
		String: em.Email,
		Valid: true,
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), nullEmail)

	if err != nil {
		log.Printf("Error creating user in database: %s", err)
		w.WriteHeader(500)
		return
    }

	mainUser := User{
		ID: dbUser.ID.UUID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		Email: dbUser.Email.String,
	}

	dat, err := json.Marshal(mainUser)
	
	if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(500)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}