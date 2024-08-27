package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	res, err := json.Marshal(data)
	if err != nil {
		statusCode = 500
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}

func responseWithErr(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}
