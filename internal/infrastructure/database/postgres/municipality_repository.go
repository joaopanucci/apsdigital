package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/repositories"
)

type MunicipalityPostgresRepository struct {
	db *pgxpool.Pool
}

func NewMunicipalityRepository(db *pgxpool.Pool) repositories.MunicipalityRepository {
	return &MunicipalityPostgresRepository{db: db}
}

func (r *MunicipalityPostgresRepository) Create(ctx context.Context, municipality *entities.Municipality) error {
	query := `
		INSERT INTO municipalities (name, ibge_code, state, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		municipality.Name,
		municipality.IBGECode,
		municipality.State,
		municipality.Active,
	).Scan(&municipality.ID)

	if err != nil {
		return fmt.Errorf("error creating municipality: %w", err)
	}

	return nil
}

func (r *MunicipalityPostgresRepository) GetByID(ctx context.Context, id int) (*entities.Municipality, error) {
	query := `
		SELECT id, name, ibge_code, state, active, created_at, updated_at
		FROM municipalities
		WHERE id = $1
	`

	municipality := &entities.Municipality{}
	var state sql.NullString
	err := r.db.QueryRow(ctx, query, id).Scan(
		&municipality.ID,
		&municipality.Name,
		&municipality.IBGECode,
		&state,
		&municipality.Active,
		&municipality.CreatedAt,
		&municipality.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("municipality not found")
		}
		return nil, fmt.Errorf("error getting municipality: %w", err)
	}

	municipality.State = state.String

	return municipality, nil
}

func (r *MunicipalityPostgresRepository) GetByName(ctx context.Context, name string) (*entities.Municipality, error) {
	query := `
		SELECT id, name, ibge_code, state, active, created_at, updated_at
		FROM municipalities
		WHERE name = $1
	`

	municipality := &entities.Municipality{}
	var state sql.NullString
	err := r.db.QueryRow(ctx, query, name).Scan(
		&municipality.ID,
		&municipality.Name,
		&municipality.IBGECode,
		&state,
		&municipality.Active,
		&municipality.CreatedAt,
		&municipality.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("municipality not found")
		}
		return nil, fmt.Errorf("error getting municipality: %w", err)
	}

	municipality.State = state.String

	return municipality, nil
}

func (r *MunicipalityPostgresRepository) List(ctx context.Context) ([]*entities.Municipality, error) {
	query := `
		SELECT id, name, ibge_code, state, active, created_at, updated_at
		FROM municipalities
		WHERE active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing municipalities: %w", err)
	}
	defer rows.Close()

	var municipalities []*entities.Municipality
	for rows.Next() {
		municipality := &entities.Municipality{}
		var state sql.NullString
		err := rows.Scan(
			&municipality.ID,
			&municipality.Name,
			&municipality.IBGECode,
			&state,
			&municipality.Active,
			&municipality.CreatedAt,
			&municipality.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning municipality: %w", err)
		}
		municipality.State = state.String
		municipalities = append(municipalities, municipality)
	}

	return municipalities, nil
}

func (r *MunicipalityPostgresRepository) Update(ctx context.Context, municipality *entities.Municipality) error {
	query := `
		UPDATE municipalities
		SET name = $2, ibge_code = $3, state = $4, active = $5, updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query,
		municipality.ID,
		municipality.Name,
		municipality.IBGECode,
		municipality.State,
		municipality.Active,
	)

	if err != nil {
		return fmt.Errorf("error updating municipality: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("municipality not found")
	}

	return nil
}

func (r *MunicipalityPostgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM municipalities WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting municipality: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("municipality not found")
	}

	return nil
}
