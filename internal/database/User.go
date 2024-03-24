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

type UserLoginResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}
