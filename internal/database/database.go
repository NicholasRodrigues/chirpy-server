package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
	NextID int           `json:"nextId"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer file.Close()
		// Initialize DBStructure with NextID set to 1
		initialStruct := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
			NextID: 1,
		}
		content, err := json.Marshal(initialStruct)
		if err != nil {
			return err
		}
		_, err = file.Write(content)
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	content, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	var dbStruct DBStructure
	err = json.Unmarshal(content, &dbStruct)
	if err != nil {
		return DBStructure{}, err
	}
	// Ensure the Chirps map is initialized
	if dbStruct.Chirps == nil {
		dbStruct.Chirps = make(map[int]Chirp)
	}
	if dbStruct.Users == nil {
		dbStruct.Users = make(map[int]User)
	}

	if dbStruct.NextID == 0 {
		dbStruct.NextID = 1
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	content, err := json.Marshal(dbStruct)
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, content, 0644)
}
