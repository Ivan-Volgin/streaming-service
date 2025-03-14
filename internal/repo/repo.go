package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"streaming-service/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type repository struct {
	pool *pgxpool.Pool
}

const (
	insertMovieQuery  = `INSERT INTO movies (uuid, title, author, description, year) VALUES ($1, $2, $3, $4, $5)`
	getAllMoviesQuery = `SELECT uuid, title, description, author, year FROM movies LIMIT $1 OFFSET $2`
	getMovieQuery     = `SELECT uuid, title, author, description, year FROM movies WHERE uuid = $1`
	updateMovieQuery  = `UPDATE movies SET title = $1, author = $2, description = $3, year = $4 WHERE uuid = $5`
	deleteMovieQuery  = `DELETE FROM movies WHERE uuid = $1`
)

type Repository interface {
	CreateMovie(ctx context.Context, film *Movie) (string, error)
	GetAllMovies(ctx context.Context, limit, offset int) (map[string]*Movie, error)
	GetMovieByID(ctx context.Context, uuid string) (*Movie, error)
	UpdateMovie(ctx context.Context, uuid string, film *Movie) error
	DeleteMovie(ctx context.Context, uuid string) error
}

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

func (r *repository) CreateMovie(ctx context.Context, movie *Movie) (string, error) {
	uuid := uuid.New().String()

	err := r.pool.QueryRow(ctx, insertMovieQuery, uuid, movie.Title, movie.Author, movie.Description, movie.Year).Scan(&uuid)
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
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan movie row")
		}
		var movie Movie
		err = rows.Scan(&movie.Title, &movie.Author, &movie.Description, &movie.Year)
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
	movie := &Movie{}

	err := r.pool.QueryRow(ctx, getMovieQuery, uuid).Scan(&movie.Title, &movie.Author, &movie.Description, &movie.Year)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query movie")
	}

	return movie, nil
}

func (r *repository) UpdateMovie(ctx context.Context, uuid string, movie *Movie) error {
	commandTag, err := r.pool.Exec(ctx, updateMovieQuery, movie.Title, movie.Description, movie.Author, movie.Year, uuid)
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
