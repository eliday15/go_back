package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with 5xx eror:", message)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v, %v", err, payload)
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}