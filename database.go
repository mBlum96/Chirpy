package main

import (
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	db := &DB{
		path: path,
		mux:  new(sync.RWMutex),
	}
	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	chirps, err := db.GetChirps()
	defer db.mux.Unlock()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:   len(chirps) + 1,
		Body: body,
	}
	chirps = append(chirps, newChirp)
	chirpMap := make(map[int]Chirp)
	for _, chirp := range chirps {
		chirpMap[chirp.Id] = chirp
	}
	err = db.writeDB(DBStructure{Chirps: chirpMap})
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error)

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error
