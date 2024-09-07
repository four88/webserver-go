package database

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string `json:"email"`
	Id           int    `json:"id"`
	Password     []byte `json: "password"`
	RefreshToken string `json: "refresh_token"`
	ExpiryTime   string `json: "expiry_time"`
	IsChirpyRed  bool   `json: "is_chirpy_red"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Id       int    `json:"id"`
	Password string `json: "password"`
}

type UserWithoutPwd struct {
	Email        string `json:"email"`
	Id           int    `json:"id"`
	RefreshToken string `json: "refresh_token"`
	IsChirpyRed  bool   `json: "is_chirpy_red"`
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
	user := User{Id: newID, Email: email, Password: hashpassword, IsChirpyRed: false}
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

			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return UserWithoutPwd{}, errors.New("Invalid Password")
			}

			refreshToken, err := GenerateRefreshToken()
			if err != nil {
				return UserWithoutPwd{}, err
			}

			dbStructure.Users[user.Id] = User{
				Email:        user.Email,
				Id:           user.Id,
				Password:     user.Password,
				RefreshToken: refreshToken,
				ExpiryTime:   time.Now().AddDate(0, 0, 60).Format("2006-01-02 15:04:05"),
				IsChirpyRed:  user.IsChirpyRed,
			}

			if err := db.writeDB(*dbStructure); err != nil {
				return UserWithoutPwd{}, err
			}

			userWithoutPwd := UserWithoutPwd{
				Email:        user.Email,
				Id:           user.Id,
				RefreshToken: refreshToken,
				IsChirpyRed:  user.IsChirpyRed,
			}

			return userWithoutPwd, nil
		}
	}

	fmt.Println("user not found 2")
	return UserWithoutPwd{}, errors.New("User not found")
}

func (db *DB) UpdateMember(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return err
	}

	dbStructure.Users[id] = User{
		Id:           user.Id,
		Email:        user.Email,
		Password:     user.Password,
		ExpiryTime:   user.ExpiryTime,
		RefreshToken: user.RefreshToken,
		IsChirpyRed:  true,
	}
	if err := db.writeDB(*dbStructure); err != nil {
		return err
	}
	return nil
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
	dbStructure.Users[id] = User{
		Id:           id,
		Email:        email,
		Password:     hashpassword,
		ExpiryTime:   user.ExpiryTime,
		RefreshToken: user.RefreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}
	if err := db.writeDB(*dbStructure); err != nil {
		return UserResponse{}, err
	}
	updatedUser := UserResponse{Id: user.Id, Email: email, Password: password}
	return updatedUser, nil
}
