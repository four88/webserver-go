package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email"`
	Id       int    `json:"id"`
	Password []byte `json: "password"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Id       int    `json:"id"`
	Password string `json: "password"`
}

type UserWithoutPwd struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	hashpassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		errors.New("Error in hashing")
		fmt.Println(err)
	}

	newID := len(dbStructure.Users) + 1
	user := User{Id: newID, Email: email, Password: hashpassword}
	dbStructure.Users[newID] = user

	// Save to disk
	if err := db.writeDB(*dbStructure); err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) Login(email string, password string) (
	UserWithoutPwd,
	error,
) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserWithoutPwd{}, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return UserWithoutPwd{}, errors.New("Invalid Password")
			}

			userWithoutPwd := UserWithoutPwd{Email: user.Email, Id: user.Id}
			return userWithoutPwd, nil
		} else {
			return UserWithoutPwd{}, errors.New("User not found")
		}
	}
	return UserWithoutPwd{}, errors.New("User not found")
}

func (db *DB) GetUserAndUpdate(id int, email string, password string) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		errors.New("User not found")
		return UserResponse{}, err
	}

	hashpassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		errors.New("Error in hashing")
		fmt.Println(err)
	}
	dbStructure.Users[id] = User{Id: id, Email: email, Password: hashpassword}
	if err := db.writeDB(*dbStructure); err != nil {
		return UserResponse{}, err
	}
	updatedUser := UserResponse{Id: user.Id, Email: email, Password: password}
	return updatedUser, nil
}
