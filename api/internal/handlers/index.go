package handlers

import (
	"github.com/tasjen/fz/db"
)

type Handlers struct {
	db *db.Queries
}

func NewHandlers(db *db.Queries) *Handlers {
	return &Handlers{db: db}
}
