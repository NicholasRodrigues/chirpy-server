package database

import (
	"fmt"
	"strconv"
)

// Create chirp without nextID

func (db *DB) CreateChirp(body string, authorId string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	intAuthorId, err := strconv.Atoi(authorId)
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		AuthorId: intAuthorId,
		Message:  body,
		ID:       id,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
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

	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		// Key does not exist, return an error indicating not found
		return Chirp{}, fmt.Errorf("chirp not found")
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStruct.Chirps[id]
	if !ok {
		return fmt.Errorf("chirp not found")
	}

	delete(dbStruct.Chirps, id)
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}
