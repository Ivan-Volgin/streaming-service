package repo

import "time"

type Movie struct {
	UUID        string    `json:"uuid"`
	OwnerID     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	Created_at  time.Time `json:"created_at"`
}

type Owner struct {
	UUID       string    `json:"uuid"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
}
