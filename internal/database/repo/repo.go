package repo

import "github.com/skamenetskiy/messages/internal/database"

func New(db database.DB) *Repo {
	return &Repo{db}
}

type Repo struct {
	db database.DB
}
