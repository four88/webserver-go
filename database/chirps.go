package database

import (
	"errors"
	"sort"
)

type Chirp struct {
	Body     string `json:"body"`
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if len(body) == 0 || len(body) > 140 {
		return Chirp{}, errors.New("invalid chirp")
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{Id: newID, Body: body, AuthorId: authorId}
	dbStructure.Chirps[newID] = chirp

	// Save to disk
	if err := db.writeDB(*dbStructure); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps(authorId int, sorting string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	if authorId != 0 {
		for _, chirp := range dbStructure.Chirps {
			if chirp.AuthorId == authorId {
				chirps = append(chirps, chirp)
			}
		}
		return chirps, nil
	} else {
		for _, chirp := range dbStructure.Chirps {
			chirps = append(chirps, chirp)
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sorting == "desc" {
			return chirps[i].Id > chirps[j].Id
		} else {
			return chirps[i].Id < chirps[j].Id
		}
	})

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chripId int, authorId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	chirp, ok := dbStructure.Chirps[chripId]
	if !ok {
		return errors.New("chirp not found")
	}
	if chirp.AuthorId == authorId {
		delete(dbStructure.Chirps, chripId)

	} else {
		return errors.New("unauthorized")
	}
	if err := db.writeDB(*dbStructure); err != nil {
		return err
	}

	return nil
}
