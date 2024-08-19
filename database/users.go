package database

type User struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newID := len(dbStructure.Users) + 1
	user := User{Id: newID, Email: email}
	dbStructure.Users[newID] = user

	// Save to disk
	if err := db.writeDB(*dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}
