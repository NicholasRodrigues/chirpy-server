package database

type Chirp struct {
	Message string `json:"body"`
	ID      int    `json:"id"`
}
