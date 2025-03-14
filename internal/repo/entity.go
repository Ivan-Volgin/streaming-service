package repo

type Movie struct {
	//UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}
