package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gbolli/chirpy/internal/auth"
)

// validate refresh token and issue new access token
func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	// check if exists
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}

	userID, err := cfg.dbQueries.GetUserFromToken(r.Context(), refreshToken)

	if err != nil {
		fmt.Print("running code for err or expired\n")
		w.WriteHeader(401)
		return
	}

	accessToken, err := auth.MakeJWT(userID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error generating token: %s\n", err)
		w.WriteHeader(500)
		return
	}

	type newToken struct {
		Token	string	`json:"token"`
	}

	token := newToken{
		Token: accessToken,
	}

	dat, err := json.Marshal(token)
	if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(500)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
