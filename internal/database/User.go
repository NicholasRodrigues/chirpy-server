package database

type User struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserResponse struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserLoginResponse struct {
	Email        string `json:"email"`
	ID           int    `json:"id"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}
