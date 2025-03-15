package repo

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx" // Алиас для migrate/pgx
	_ "github.com/golang-migrate/migrate/v4/source/file"           // Поддержка чтения миграций из файлов
	"github.com/jackc/pgx/v5"                                      // Основной пакет pgx
	"github.com/jackc/pgx/v5/pgxpool"                              // Пул соединений
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"streaming-service/internal/config"
)

type repository struct {
	pool *pgxpool.Pool
}

type Repositories struct {
	MovieRepo MovieRepository
	OwnerRepo OwnerRepository
}

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (*Repositories, error) {
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

	if err := applyMigrations(pool, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to apply migrations")
	}
	
	baseRepo := &repository{pool}

	return &Repositories{
		MovieRepo: baseRepo,
		OwnerRepo: baseRepo,
	}, nil
}

func applyMigrations(pool *pgxpool.Pool, cfg config.PostgreSQL) error {
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

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}
