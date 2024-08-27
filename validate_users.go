package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/four88/webserver/database"
)

type requestUserPostBody struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ExipiresIn int    `json:"expires_in_seconds"`
}

type responsePostUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type responsePutUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Id       int    `json:"id"`
}

type requestUserPutBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseLogin struct {
	responsePostUser
	Token string `json:"token"`
}

func createUser(w http.ResponseWriter, r *http.Request, db database.DB) {
	var req requestUserPostBody
	var res responsePostUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	statusCode := 201
	if err != nil {
		msg := "Error decoding JSON"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}
	data, err := db.CreateUser(req.Email, req.Password)
	if err != nil {
		msg := "Error creating user"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}

	res.Email = data.Email
	res.Id = data.Id

	responseWithJSON(w, res, statusCode)
}

func login(w http.ResponseWriter, r *http.Request, db database.DB, secretKey string) {
	var req requestUserPostBody
	var res responseLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	statusCode := 200
	if err != nil {
		msg := "Error decoding JSON"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}
	data, err := db.Login(req.Email, req.Password)
	if err != nil {
		msg := "Unauthorized"
		statusCode = 401
		responseWithErr(w, msg, statusCode)
		return
	}
	if req.ExipiresIn == 0 {
		req.ExipiresIn = 86400
	}
	// gerate token
	token, err := createJWT(strconv.Itoa(data.Id), secretKey, int64(req.ExipiresIn))
	if err != nil {
		fmt.Println(err)
		msg := "Error generating token"
		statusCode = 500
		responseWithErr(w, msg, statusCode)
		return
	}
	res.Email = data.Email
	res.Id = data.Id
	res.Token = token
	responseWithJSON(w, res, statusCode)
}

func handleCheckToken(w http.ResponseWriter, r *http.Request, db database.DB, secretKey string) {
	var req requestUserPutBody
	var res responsePutUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	statusCode := 200
	if err != nil {
		msg := "Error decoding JSON"
		statusCode = 401
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
	userID, err := checkAndClaimToken(tokenString, secretKey)
	if err != nil {
		msg := "Invalid token"
		responseWithErr(w, msg, 401)
		return
	}
	updatedUser, err := db.GetUserAndUpdate(userID, req.Email, req.Password)
	if err != nil {
		msg := "Not found this user"
		responseWithErr(w, msg, 401)
		return
	}
	res = responsePutUser{Id: updatedUser.Id, Password: updatedUser.Password, Email: updatedUser.Email}
	responseWithJSON(w, res, statusCode)
}
