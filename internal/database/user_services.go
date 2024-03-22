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
		ID:       id,
		Email:    email,
		Password: string(encryptedPassword),
	}
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
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

func (db *DB) LoginUser(email, password string) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return UserResponse{}, fmt.Errorf("invalid password")
			}

			return UserResponse{
				ID:    user.ID,
				Email: user.Email,
			}, nil
		}
	}

	return UserResponse{}, fmt.Errorf("user not found")
}
