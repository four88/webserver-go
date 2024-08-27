package database

import (
	"errors"
	"sort"
)

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if len(body) == 0 || len(body) > 140 {
		return Chirp{}, errors.New("invalid chirp")
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{Id: newID, Body: body}
	dbStructure.Chirps[newID] = chirp

	// Save to disk
	if err := db.writeDB(*dbStructure); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
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
