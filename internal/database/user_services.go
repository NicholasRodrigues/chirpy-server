package database

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string) (UserResponse, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}

	id := len(dbStruct.Users) + 1
	user := User{
		ID:          id,
		Email:       email,
		Password:    string(encryptedPassword),
		IsChirpyRed: false,
	}
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		Email:       user.Email,
		ID:          user.ID,
		IsChirpyRed: user.IsChirpyRed,
	}, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, fmt.Errorf("chirp not found")
	}

	return user, nil
}

func (db *DB) LoginUser(email, password string) (UserLoginResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserLoginResponse{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return UserLoginResponse{}, fmt.Errorf("invalid password")
			}

			return UserLoginResponse{
				Email:       user.Email,
				ID:          user.ID,
				IsChirpyRed: user.IsChirpyRed,
			}, nil
		}
	}

	return UserLoginResponse{}, fmt.Errorf("user not found")
}

func (db *DB) UpdateUser(id int, email, password string) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return UserResponse{}, fmt.Errorf("user not found")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}

	user.Email = email
	user.Password = string(encryptedPassword)
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		Email:       user.Email,
		ID:          user.ID,
		IsChirpyRed: user.IsChirpyRed,
	}, nil
}

func (db *DB) UpdateUserChirpyRed(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}

	user.IsChirpyRed = true
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
