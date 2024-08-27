package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

func generateRefreshToken() (string, error) {
	// Generate 256-bit random token
	c := 32
	randomByte := make([]byte, c)
	_, err := rand.Read(randomByte)
	if err != nil {
		log.Printf("Error generating random token: %v", err)
		return "", err
	}

	randomToken := hex.EncodeToString(randomByte)

	return randomToken, nil
}

func (db *DB) ValidateRefreshToken(refreshToken string) (string, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return "", err
	}
	for _, user := range dbStruct.Users {
		if user.RefreshToken == refreshToken {

			// Log the ExpiryTime to see its format
			log.Printf("User.ExpiryTime: %s", user.ExpiryTime)

			// Define the correct layout based on the actual format of ExpiryTime
			layout := "2006-01-02 15:04:05"

			// Parse the string to time.Time
			parsedTime, err := time.Parse(layout, user.ExpiryTime)
			if err != nil {
				log.Fatalf("Error parsing time: %v", err)
				return "", errors.New("Error parsing time")
			}

			if time.Now().After(parsedTime) {
				log.Println("Token expired")
				return "", errors.New("Token expired")
			}
			return user.RefreshToken, nil
		}
	}
	log.Println("Refresh token not found")
	return "", errors.New("Refresh token not found")
}
