package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/pkg/errors"
)

const (
	insertMovieQuery  = `INSERT INTO movies (uuid, owner_id, title, author, description, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING uuid`
	getAllMoviesQuery = `SELECT uuid, title, description, author, year FROM movies LIMIT $1 OFFSET $2`
	getMovieQuery     = `SELECT title, author, description, year FROM movies WHERE uuid = $1`
	updateMovieQuery  = `UPDATE movies SET title = $1, author = $2, description = $3, year = $4 WHERE uuid = $5`
	deleteMovieQuery  = `DELETE FROM movies WHERE uuid = $1`
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *Movie, ownerName string) (string, error)
	GetAllMovies(ctx context.Context, limit, offset int) (map[string]*Movie, error)
	GetMovieByID(ctx context.Context, uuid string) (*Movie, error)
	UpdateMovie(ctx context.Context, uuid string, film *Movie) error
	DeleteMovie(ctx context.Context, uuid string) error
}

func (r *repository) CreateMovie(ctx context.Context, movie *Movie, ownerName string) (string, error) {
	uuid := uuid.New().String()

	var ownerUUID string
	owner, err := r.GetOwnerByName(ctx, ownerName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			newOwnerUUID, err := r.CreateOwner(ctx, &Owner{Name: ownerName})
			if err != nil {
				return "", fmt.Errorf("owner of movie not found, failed to create new owner: %w", err)
			}
			ownerUUID = newOwnerUUID
		} else {
			// Другая ошибка при выполнении запроса
			return "", fmt.Errorf("owner of movie not found, failed to query owner: %w", err)
		}
	} else {
		ownerUUID = owner.UUID
	}

	err = r.pool.QueryRow(ctx, insertMovieQuery, uuid, ownerUUID, movie.Title, movie.Author, movie.Description, movie.Year).Scan(&uuid)
	if err != nil {
		return "", errors.Wrap(err, "failed to insert movie")
	}
	return uuid, nil
}

func (r *repository) GetAllMovies(ctx context.Context, limit, offset int) (map[string]*Movie, error) {
	movies := make(map[string]*Movie)

	rows, err := r.pool.Query(ctx, getAllMoviesQuery, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query all movies")
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		var uuid string

		err := rows.Scan(&uuid, &movie.Title, &movie.Author, &movie.Description, &movie.Year)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan movie row")
		}
		movies[uuid] = &movie
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error occurred during iteration over movie rows")
	}

	return movies, nil
}

func (r *repository) GetMovieByID(ctx context.Context, uuid string) (*Movie, error) {
	movie := &Movie{UUID: uuid}

	err := r.pool.QueryRow(ctx, getMovieQuery, uuid).Scan(&movie.Title, &movie.Author, &movie.Description, &movie.Year)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query movie")
	}

	return movie, nil
}

func (r *repository) UpdateMovie(ctx context.Context, uuid string, movie *Movie) error {
	commandTag, err := r.pool.Exec(ctx, updateMovieQuery, movie.Title, movie.Author, movie.Description, movie.Year, uuid)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows updated, movie with given UUID not found")
	}

	return nil
}

func (r *repository) DeleteMovie(ctx context.Context, uuid string) error {
	commandTag, err := r.pool.Exec(ctx, deleteMovieQuery, uuid)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete query")
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows deleted, movie with given UUID not found")
	}

	return nil
}
