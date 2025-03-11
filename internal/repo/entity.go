package repo

type Movie struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}
