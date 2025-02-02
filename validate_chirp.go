package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirp(writer http.ResponseWriter, req *http.Request) {

	type parameters struct {
        Body string `json:"body"`
    }

	type validVal struct {
        CleanedBody string `json:"cleaned_body"`
    }

	type errorVal struct {
		Errormsg string `json:"error"`
	}

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
    }

	writer.Header().Set("Content-Type", "application/json")

	// Test chirp logic
	if len(params.Body) > 140 {
		respBody := errorVal {
			Errormsg: "Chirp is too long",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(400)
    	writer.Write(dat)
		return
	}

	respBody := validVal {
		CleanedBody: cleanBody(params.Body),
	}

    dat, err := json.Marshal(respBody)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			writer.WriteHeader(500)
			return
	}
    
    writer.WriteHeader(200)
    writer.Write(dat)
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