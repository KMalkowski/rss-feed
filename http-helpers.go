package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln("Error marshaling json body.")
		return
	}

	w.WriteHeader(code)
	w.Write(body)
	return
}

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	type HealthzResponseBody struct {
		Status string `json:"status"`
	}
	respondWithJson(w, http.StatusOK, HealthzResponseBody{Status: "ok"})
	return
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, ErrorResponse{Error: message})
	return
}

func ErrHandler(w http.ResponseWriter, req *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	return
}
