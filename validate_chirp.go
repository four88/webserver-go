package main

import (
	"encoding/json"
	"net/http"

	"github.com/four88/webserver/database"
)

type requestBody struct {
	Body string `json:"body"`
}
type responsePostBody struct {
	Email string `json:"body"`
	Id    int    `json:"id"`
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
