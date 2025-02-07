package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)


func (cfg *apiConfig) updateUserToRed(w http.ResponseWriter, r *http.Request) {

	// structs

	type parameters struct {
		Event string `json:"event"`
        Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
    }

	// decode params

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	// ignore anything other than "user.upgraded"
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	// update user in database

	err = cfg.dbQueries.UpgradeUserToRed(r.Context(), params.Data.UserID)
	if err != nil {
		log.Printf("Error updating user in database: %s", err)
		w.WriteHeader(404)
		return
    }

	// success

	w.WriteHeader(204)
}

