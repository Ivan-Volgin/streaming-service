package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type repository struct {
	mu     sync.RWMutex
	movies map[uuid.UUID]*Movie
}

type Repository interface {
	CreateMovie(ctx context.Context, film *Movie) (uuid.UUID, error)
	GetAllMovies(ctx context.Context) (map[uuid.UUID]*Movie, error)
	GetMovieByID(ctx context.Context, uuid uuid.UUID) (*Movie, error)
	UpdateMovie(ctx context.Context, film *Movie) error
	DeleteMovie(ctx context.Context, uuid uuid.UUID) error
}

func NewRepository(ctx context.Context) (Repository, error) {
	return &repository{movies: make(map[uuid.UUID]*Movie)}, nil
}

func (r *repository) CreateMovie(ctx context.Context, movie *Movie) (uuid.UUID, error) {
	if err := r.checkIfExists(movie); err != nil {
		return uuid.Nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	uuid := uuid.New()
	movie.UUID = uuid
	r.movies[uuid] = movie

	return uuid, nil
}

func (r *repository) GetAllMovies(ctx context.Context) (map[uuid.UUID]*Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.movies, nil
}

func (r *repository) GetMovieByID(ctx context.Context, uuid uuid.UUID) (*Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	movie, exists := r.movies[uuid]
	if !exists {
		return nil, fmt.Errorf("Film with uuid %s not found", uuid)
	}

	return movie, nil
}

func (r *repository) UpdateMovie(ctx context.Context, movie *Movie) error {
	r.mu.RLock()
	_, exists := r.movies[movie.UUID]
	if !exists {
		return fmt.Errorf("Film with uuid %s not found", movie.UUID)
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	r.movies[movie.UUID] = movie // так как в задании сказано использовать метод PUT, я просто новую версию фильма (можно было бы обновить конкретные поля, но тогда должен быть метод PATCH)
	return nil
}

func (r *repository) DeleteMovie(ctx context.Context, uuid uuid.UUID) error {
	r.mu.RLock()
	_, exists := r.movies[uuid]
	if !exists {
		return fmt.Errorf("Film with uuid %s not found", uuid)
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.movies, uuid)
	
	return nil
}

func (r *repository) checkIfExists(movie *Movie) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, existingFilm := range r.movies {
		if existingFilm.Title == movie.Title && existingFilm.Author == movie.Author {
			return fmt.Errorf("Film with title %s already exists", movie.Title)
		}
	}
	return nil
}
