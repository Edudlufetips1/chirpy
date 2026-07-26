package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	type chirpError struct {
		Error string `json:"error"`
	}

	type chirpRequest struct {
		Body string `json:"body"`
	}

	type chirpCleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	requests := chirpRequest{}
	err := decoder.Decode(&requests)
	if err != nil {
		log.Printf("Error decoding validation: %s", err)
		w.WriteHeader(400)
		return
	}

	if len(requests.Body) > 140 {
		response := chirpError{
			Error: "Chirp body exceeds 140 characters",
		}
		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error encoding validation response: %s", err)
			w.WriteHeader(400)
			return
		}
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	filteredChirp := filterProfanity(requests.Body)

	response := chirpCleanedResponse{
		CleanedBody: filteredChirp,
	}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(200)
	w.Write(jsonResponse)
}
