package service

type CreateMovieRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}
type UpdateMovieRequest struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

type GetMovieRequest struct {
	UUID string `json:"uuid"`
}

type DeleteMovieRequest struct {
	UUID string `json:"uuid"`
}
