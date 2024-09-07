package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/four88/webserver/database"
)

type requestBody struct {
	Body string `json:"body"`
}
type responsePostBody struct {
	Email    string `json:"body"`
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
}

func createChirp(w http.ResponseWriter, r *http.Request, db database.DB, secretKey string) {

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

	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if len(tokenString) <= 0 {
		msg := "No token provided"
		responseWithErr(w, msg, 401)
		return
	}
	authorId, err := checkAndClaimToken(tokenString, secretKey)
	if err != nil {
		msg := "Invalid token"
		statusCode := 401
		responseWithErr(w, msg, statusCode)
		return
	}

	if len(req.Body) > 140 {
		msg := "Chirp is too long"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	} else {
		data, err := db.CreateChirp(req.Body, authorId)
		if err != nil {
			msg := "Error creating chirp"
			statusCode = 400
			responseWithErr(w, msg, statusCode)
		}

		res.Email = data.Body
		res.Id = data.Id
		res.AuthorId = data.AuthorId
	}

	responseWithJSON(w, res, statusCode)
}

func getChirps(w http.ResponseWriter, r *http.Request, db database.DB, secretKey string) {
	authorId := 0
	query := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")
	if query != "" {
		idString, err := strconv.Atoi(query)
		if err != nil {
			fmt.Println(err)
		}

		authorId = idString
	}

	chirps, err := db.GetChirps(authorId, sorting)
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

func deleteChrips(w http.ResponseWriter, r *http.Request, db database.DB, id int, secretKey string) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if len(tokenString) <= 0 {
		msg := "No token provided"
		responseWithErr(w, msg, 401)
		return
	}
	userId, err := checkAndClaimToken(tokenString, secretKey)
	if err != nil {
		msg := "Invalid token"
		responseWithErr(w, msg, 401)
		return
	}

	err = db.DeleteChirp(id, userId)
	if err != nil {
		msg := "User not authorized to delete this chirp"
		responseWithErr(w, msg, 403)
	}

	responseWithJSON(w, "", 204)
}
