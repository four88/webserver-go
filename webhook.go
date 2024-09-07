package main

import (
	"encoding/json"
	"github.com/four88/webserver/database"
	"net/http"
	"strings"
)

type requestWebHook struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func updateMemberHook(w http.ResponseWriter, r *http.Request, db database.DB, apiKey string) {
	apiHeader := r.Header.Get("Authorization")
	key := strings.TrimPrefix(apiHeader, "ApiKey ")
	if key != apiKey {
		responseWithErr(w, "Invalid API Key", 401)
	}
	var req requestWebHook
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	statusCode := 204
	if err != nil {
		msg := "Error decoding JSON"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}

	if req.Event != "user.upgraded" {
		responseWithErr(w, "", 204)
		return
	}
	err = db.UpdateMember(req.Data.UserID)
	if err != nil {
		msg := "Error updating user"
		statusCode = 400
		responseWithErr(w, msg, statusCode)
		return
	}

	responseWithJSON(w, "", statusCode)
}
