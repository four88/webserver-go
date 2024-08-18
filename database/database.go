package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)



type DB struct {
  path string
  mux *sync.RWMutex
}


type DBStructure struct {
  Chirps map[int]Chirp  `json:"chirps"`
}

type Chirp struct {
  Body string `json:"body"`
  Id int `json:"id"`
}

func NewDB(path string) (*DB, error) {
   db := &DB{
        path: path,
        mux:  &sync.RWMutex{},
    }
    if err := db.ensureDB(); err != nil {
        return nil, err
    }
    return db, nil
}


func (db *DB) ensureDB() error {
    if _, err := os.Stat(db.path); os.IsNotExist(err) {
        initialData := DBStructure{
            Chirps: make(map[int]Chirp),
        }
        return db.writeDB(initialData)
    }
    return nil
}


func (db *DB) loadDB() (*DBStructure, error) {
  db.mux.Lock()
  defer db.mux.Unlock()
  f, err := os.Open(db.path)
  if err != nil {
    return nil, err
  }
  defer f.Close()
  var dbStruct DBStructure
  if err := json.NewDecoder(f).Decode(&dbStruct); err != nil {
    return nil, err
  }
  return &dbStruct, nil
}


func (db *DB) writeDB(dbStructure DBStructure) error {
    db.mux.Lock()
    defer db.mux.Unlock()

    data, err := json.Marshal(dbStructure)
    if err != nil {
        return err
    }

    return os.WriteFile(db.path, data, 0644)
}

func (db *DB) CreateChirp(body string)(Chirp, error){
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
