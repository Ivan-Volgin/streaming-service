package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type in_memory_repository struct {
	mu     sync.RWMutex
	movies map[string]*Movie
}

func newRepository(ctx context.Context) (Repository, error) {
	return &in_memory_repository{movies: make(map[string]*Movie)}, nil
}

func (r *in_memory_repository) CreateMovie(ctx context.Context, movie *Movie) (string, error) {
	if err := r.checkIfExists(movie); err != nil {
		return "", err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	uuid := uuid.New().String()
	r.movies[uuid] = movie

	return uuid, nil
}

func (r *in_memory_repository) GetAllMovies(ctx context.Context, limit, offset int) (map[string]*Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.movies, nil
}

func (r *in_memory_repository) GetMovieByID(ctx context.Context, uuid string) (*Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	movie, exists := r.movies[uuid]
	if !exists {
		return nil, fmt.Errorf("Film with uuid %s not found", uuid)
	}

	return movie, nil
}

func (r *in_memory_repository) UpdateMovie(ctx context.Context, uuid string, movie *Movie) error {
	r.mu.RLock()
	_, exists := r.movies[uuid]
	if !exists {
		return fmt.Errorf("Film with uuid %s not found", uuid)
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	r.movies[uuid] = movie // так как в задании сказано использовать метод PUT, я просто новую версию фильма (можно было бы обновить конкретные поля, но тогда должен быть метод PATCH)
	return nil
}

func (r *in_memory_repository) DeleteMovie(ctx context.Context, uuid string) error {
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

func (r *in_memory_repository) checkIfExists(movie *Movie) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, existingFilm := range r.movies {
		if existingFilm.Title == movie.Title && existingFilm.Author == movie.Author {
			return fmt.Errorf("Film with title %s already exists", movie.Title)
		}
	}
	return nil
}
