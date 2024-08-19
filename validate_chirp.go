package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/four88/webserver/database"
)

type requestBody struct {
	Body string `json:"body"`
}
type responsePostBody struct {
	Email string `json:"body"`
	Id    int    `json:"id"`
}

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

func profaneWords(sentence string) string {
	profaneWords := []string{"sharbert", "kerfuffle", "fornax"}
	listWord := strings.Split(sentence, " ")

	for i, word := range listWord {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				listWord[i] = "****"
			}
		}
	}

	return strings.Join(listWord, " ")
}

func createChirp(w http.ResponseWriter, r *http.Request, db database.DB) {
	var req requestBody
	var res responsePostBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	statusCode := 201
	if err != nil {
		msg := "Error decoding JSON"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}
	if len(req.Body) > 140 {
		msg := "Chirp is too long"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	} else {
		data, err := db.CreateChirp(req.Body)
		if err != nil {
			msg := "Error creating chirp"
			statusCode = 400
			responseWithErr(w, msg, statusCode)
		}

		res.Email = data.Body
		res.Id = data.Id
	}

	responseWithJSON(w, res, statusCode)
}

func getChirps(w http.ResponseWriter, r *http.Request, db database.DB) {
	chirps, err := db.GetChirps()
	if err != nil {
		msg := "Error getting chirps"
		responseWithErr(w, msg, 500)
		return
	}
	responseWithJSON(w, chirps, 200)
}

func getChirp(w http.ResponseWriter, r *http.Request, db database.DB, id int) {
	chirp, err := db.GetChirp(id)
	if err != nil {
		msg := "Error getting chirp"
		responseWithErr(w, msg, 404)
		return
	}
	responseWithJSON(w, chirp, 200)
}
