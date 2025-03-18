package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	insertOwnerQuery    = `INSERT INTO owners (uuid, name) VALUES ($1, $2) RETURNING uuid`
	getAllOwnersQuery   = `SELECT uuid, name, created_at FROM owners LIMIT $1 OFFSET $2`
	getOwnerByIdQuery   = `SELECT name, created_at FROM owners WHERE uuid = $1`
	getOwnerByNameQuery = `SELECT uuid, created_at FROM owners WHERE name = $1`
	updateOwnerQuery    = `UPDATE owners SET name = $1 WHERE uuid = $2`
	deleteOwnerQuery    = `DELETE FROM owners WHERE uuid = $1`
)

type OwnerRepository interface {
	CreateOwner(ctx context.Context, owner *Owner) (string, error)
	GetAllOwners(ctx context.Context, limit, offset int) (map[string]*Owner, error)
	GetOwnerByID(ctx context.Context, uuid string) (*Owner, error)
	GetOwnerByName(ctx context.Context, name string) (*Owner, error)
	UpdateOwner(ctx context.Context, uuid string, film *Owner) error
	DeleteOwner(ctx context.Context, uuid string) error
}

func (r *repository) CreateOwner(ctx context.Context, owner *Owner) (string, error) {
	uuid := uuid.New().String()

	err := r.pool.QueryRow(ctx, insertOwnerQuery, uuid, owner.Name).Scan(&uuid)
	if err != nil {
		return "", errors.Wrap(err, "failed to insert owner")
	}
	return uuid, nil
}

func (r *repository) GetAllOwners(ctx context.Context, limit, offset int) (map[string]*Owner, error) {
	owners := make(map[string]*Owner)

	rows, err := r.pool.Query(ctx, getAllOwnersQuery, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query all owners")
	}
	defer rows.Close()

	for rows.Next() {
		var owner Owner

		err = rows.Scan(&owner.UUID, &owner.Name, &owner.Created_at)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan owner row")
		}
		owners[owner.UUID] = &owner
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error occurred during iteration over owner rows")
	}

	return owners, nil
}

func (r *repository) GetOwnerByID(ctx context.Context, uuid string) (*Owner, error) {
	owner := &Owner{UUID: uuid}

	err := r.pool.QueryRow(ctx, getOwnerByIdQuery, uuid).Scan(&owner.Name, &owner.Created_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, "owner not found")
		}
		return nil, errors.Wrap(err, "failed to query owner by uuid")
	}

	return owner, nil
}

func (r *repository) GetOwnerByName(ctx context.Context, name string) (*Owner, error) {
	owner := &Owner{Name: name}

	err := r.pool.QueryRow(ctx, getOwnerByNameQuery, name).Scan(&owner.UUID, &owner.Created_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, "owner not found")
		}
		return nil, errors.Wrap(err, "failed to query owner by uuid")
	}

	return owner, nil
}

func (r *repository) UpdateOwner(ctx context.Context, uuid string, owner *Owner) error {
	commandTag, err := r.pool.Exec(ctx, updateOwnerQuery, owner.Name, uuid)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows updated, owner with given UUID not found")
	}

	return nil
}

func (r *repository) DeleteOwner(ctx context.Context, uuid string) error {
	commandTag, err := r.pool.Exec(ctx, deleteOwnerQuery, uuid)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete query")
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows deleted, owner with given UUID not found")
	}

	return nil
}
