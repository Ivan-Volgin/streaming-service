package repo

import (
	"github.com/google/uuid"
)

type Movie struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
}
