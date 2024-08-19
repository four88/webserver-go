package main

import (
	"encoding/json"
	"github.com/four88/webserver/database"
	"net/http"
)

type requestUserPostBody struct {
	Email string `json:"email"`
}

type respoonsePostUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request, db database.DB) {
	var req requestUserPostBody
	var res respoonsePostUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	statusCode := 201
	if err != nil {
		msg := "Error decoding JSON"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}
	data, err := db.CreateUser(req.Email)
	if err != nil {
		msg := "Error creating chirp"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
	}

	res.Email = data.Email
	res.Id = data.Id

	responseWithJSON(w, res, statusCode)
}
