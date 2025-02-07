package main

import (
	"log"
	"net/http"

	"github.com/gbolli/chirpy/internal/auth"
)

// revoke token by adding revoked_at time
func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	// check if exists
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}

	err = cfg.dbQueries.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error revoking token: %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)

}
