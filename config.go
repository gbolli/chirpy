package main

import (
	"sync/atomic"

	"github.com/gbolli/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
	secret string
}
