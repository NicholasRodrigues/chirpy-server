package database

import (
	"fmt"
	"sort"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		ID:      dbStruct.NextID, // Assign current NextID to chirp
		Message: body,
	}
	dbStruct.Chirps[chirp.ID] = chirp
	dbStruct.NextID++ // Increment for next chirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var ids []int
	for id := range dbStruct.Chirps {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, id := range ids {
		chirps = append(chirps, dbStruct.Chirps[id])
	}

	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

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
