package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/four88/webserver/database"
)

type responseRefreshToken struct {
	Token string `json:"token"`
}

func handleRefresh(w http.ResponseWriter, r *http.Request, db database.DB, secretKey string) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	statusCode := 200
	if len(tokenString) <= 0 {
		msg := "No token provided"
		statusCode = 401
		responseWithErr(w, msg, statusCode)
		return
	}
	userID, err := db.ValidateRefreshToken(tokenString)
	if err != nil {
		msg := "Unauthorized"
		statusCode = 401
		responseWithErr(w, msg, statusCode)
		return
	}
	newJWT, err := createJWT(strconv.Itoa(userID), secretKey)
	if err != nil {
		msg := "Error generating token"
		statusCode = 401
		responseWithErr(w, msg, statusCode)
	}
	res := responseRefreshToken{Token: newJWT}
	responseWithJSON(w, res, statusCode)
}
func handleRevolk(w http.ResponseWriter, r *http.Request, db database.DB) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	statusCode := 204
	if len(tokenString) <= 0 {
		msg := "No token provided"
		statusCode = 401
		responseWithErr(w, msg, statusCode)
		return
	}
	_, err := db.RevolkRefreshToken(tokenString)
	if err != nil {
		msg := "Unauthorized"
		statusCode = 401
		responseWithErr(w, msg, statusCode)
		return
	}

	res := ""
	responseWithJSON(w, res, statusCode)
}
