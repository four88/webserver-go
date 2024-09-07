package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

func GenerateRefreshToken() (string, error) {
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

func (db *DB) RevolkRefreshToken(refreshToken string) (string, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return "", err
	}
	userData := User{}
	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			userData = user
		} else {
			return "", errors.New("Refresh token not found")
		}
	}

	newToken, err := GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	dbStructure.Users[userData.Id] = User{
		Email:        userData.Email,
		Id:           userData.Id,
		Password:     userData.Password,
		RefreshToken: newToken,
		ExpiryTime:   time.Now().Add(time.Duration(60) * time.Second).Format("2006-01-02 15:04:05"),
	}
	if err := db.writeDB(*dbStructure); err != nil {
		return "", err
	}
	return "", nil
}

func (db *DB) ValidateRefreshToken(refreshToken string) (int, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return 0, err
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
				return 0, errors.New("Error parsing time")
			}

			if time.Now().After(parsedTime) {
				log.Println("Token expired")
				return 0, errors.New("Token expired")
			}
			return user.Id, nil
		}
	}
	log.Println("Refresh token not found")
	return 0, errors.New("Refresh token not found")
}

func (db *DB) GetUserByRefreshToken(refreshToken string) (int, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return 0, err
	}
	for _, user := range dbStruct.Users {
		if user.RefreshToken == refreshToken {
			return user.Id, nil
		}
	}
	return 0, errors.New("Refresh token not found")
}
