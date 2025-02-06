package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gbolli/chirpy/internal/auth"
	"github.com/google/uuid"
)


func (cfg *apiConfig) userLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
        Email string `json:"email"`
		ExpiresInSeconds *int `json:"expires_in_seconds"`
    }

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token	  string	`json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	// check password hash

	nullEmail := sql.NullString{
		String: params.Email,
		Valid: true,
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), nullEmail)
	if err != nil || auth.CheckPasswordHash(params.Password, dbUser.HashedPassword) != nil {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
		return
    }

	// create token
	expiresInSeconds := 3600
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds < 3600 { expiresInSeconds = *params.ExpiresInSeconds}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Duration(expiresInSeconds * int(time.Second)))
	if err != nil {
		log.Printf("Error generating token: %s\n", err)
		w.WriteHeader(500)
		return
	}

	mainUser := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		Email: dbUser.Email.String,
		Token: token,
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