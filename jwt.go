package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(userID string, secretKey string, expiresInSeconds int64) (string, error) {
	// Create the JWT claims
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   userID,
	}

	// Create the token using the HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	secretKeyBit := []byte(secretKey)
	tokenString, err := token.SignedString(secretKeyBit)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func checkAndClaimToken(tokenString string, jwtSecret string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if claim, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userIdString := claim.Subject
		userID, err := strconv.Atoi(userIdString)
		if err != nil {
			return 0, err
		}
		return userID, nil
	} else {
		return 0, err
	}
}
