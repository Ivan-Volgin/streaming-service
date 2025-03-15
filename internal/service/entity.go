package service

type CreateMovieRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	OwnerName   string `json:"owner_name"`
}

type CreateOwnerRequest struct {
	Name string `json:"name"`
}
type UpdateMovieRequest struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

type UpdateOwnerRequest struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
type GetMovieRequest struct {
	UUID string `json:"uuid"`
}

type GetOwnerByUUIDRequest struct {
	UUID string `json:"uuid"`
}

type GetOwnerByNameRequest struct {
	Name string `json:"name"`
}

type DeleteMovieRequest struct {
	UUID string `json:"uuid"`
}

type DeleteOwnerRequest struct {
	UUID string `json:"uuid"`
}
