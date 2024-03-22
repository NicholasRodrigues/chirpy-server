package database

type User struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}
