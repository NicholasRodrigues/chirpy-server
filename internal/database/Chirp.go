package database

type Chirp struct {
	AuthorId int    `json:"author_id"`
	Message  string `json:"body"`
	ID       int    `json:"id"`
}
