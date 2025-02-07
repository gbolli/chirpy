package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gbolli/chirpy/internal/auth"
	"github.com/gbolli/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

	// structs

	type parameters struct {
		Password string `json:"password"`
        Email string `json:"email"`
    }

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		IsChirpyRed	bool	`json:"is_chirpy_red"`
	}

	// authenticate user token

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

	// decode params

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)

    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	nullEmail := sql.NullString{
		String: params.Email,
		Valid: true,
	}

	// hash password

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
    }

	// set database params

	updateParams := database.UpdateUserParams{
		ID: tokenUserID,
		HashedPassword: hashed_password,
		Email: nullEmail,
	}

	// update user in database

	dbUser, err := cfg.dbQueries.UpdateUser(r.Context(), updateParams)
	if err != nil {
		log.Printf("Error updating user in database: %s", err)
		w.WriteHeader(500)
		return
    }

	// prepare reply data

	mainUser := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		Email: dbUser.Email.String,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	dat, err := json.Marshal(mainUser)
	
	if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(500)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

