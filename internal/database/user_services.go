package database

import "sort"

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:    dbStruct.NextID,
		Email: email,
	}

	dbStruct.Users[user.ID] = user
	dbStruct.NextID++

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var ids []int
	for id := range dbStruct.Users {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	users := make([]User, 0, len(dbStruct.Users))
	for _, id := range ids {
		users = append(users, dbStruct.Users[id])
	}

	return users, nil
}
