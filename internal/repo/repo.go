package repo

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"streaming-service/internal/config"
	"time"
)

type repository struct {
	pool *pgxpool.Pool
}

type Repositories interface {
	MovieRepository
	OwnerRepository
	Close()
}

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repositories, error) {
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

	if err := applyMigrations(pool); err != nil {
		return nil, errors.Wrap(err, "failed to apply migrations")
	}

	return &repository{pool}, nil
}

func applyMigrations(pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	driver, err := pgxMigrate.WithInstance(sqlDB, &pgxMigrate.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to initialize pgx migrate driver")
	}

	migrationsPath := "./migrations"
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}

func (r *repository) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем запрос, который задерживает выполнение на 3 секунды, чтобы увидеть 1 используемый коннект
	go func() {
		fmt.Println("Выполняется запрос SELECT pg_sleep(3)...")
		_, err := r.pool.Exec(ctx, "SELECT pg_sleep(3)")
		if err != nil {
			fmt.Errorf("Ошибка выполнения запроса: %v", err)
		}

		fmt.Println("Запрос выполнен успешно")
	}()

	time.Sleep(5 * time.Second)
	stats := r.pool.Stat()

	fmt.Println(float64(stats.MaxConns()))
	fmt.Println(float64(stats.TotalConns())) // занято 1 соединение

	if r.pool != nil {
		r.pool.Close()
	}

	stats = r.pool.Stat()

	fmt.Println(float64(stats.MaxConns()))
	fmt.Println(float64(stats.TotalConns())) // 0 занятых соединений
}
