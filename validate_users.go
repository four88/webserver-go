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
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
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
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
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
	res.IsChirpyRed = data.IsChirpyRed

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
		req.ExipiresIn = 3600
	}
	// gerate token
	token, err := createJWT(strconv.Itoa(data.Id), secretKey)
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
	res.RefreshToken = data.RefreshToken
	res.IsChirpyRed = data.IsChirpyRed
	responseWithJSON(w, res, statusCode)
}

// TODO : need to check refresh toekn and get user id
// now worked only jwt token. Refresh token need to be checked
func handleUpdateUser(w http.ResponseWriter, r *http.Request, db database.DB, secretKey string) {
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
