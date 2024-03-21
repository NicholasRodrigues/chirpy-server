package database

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

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
