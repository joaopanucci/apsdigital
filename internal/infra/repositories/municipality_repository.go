package repositories

import (
	"context"
	"fmt"

	"apsdigital/internal/domain/entities"
	"apsdigital/internal/infra/db"

	"github.com/google/uuid"
)

type municipalityRepository struct {
	db *db.PostgresDB
}

func NewMunicipalityRepository(db *db.PostgresDB) *municipalityRepository {
	return &municipalityRepository{db: db}
}

func (r *municipalityRepository) Create(ctx context.Context, municipality *entities.Municipality) error {
	query := `
		INSERT INTO municipalities (name, state, active, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	
	err := r.db.Pool.QueryRow(ctx, query, 
		municipality.Name, 
		municipality.State, 
		municipality.Active,
	).Scan(&municipality.ID)
	
	if err != nil {
		return fmt.Errorf("error creating municipality: %w", err)
	}
	
	return nil
}

func (r *municipalityRepository) GetByID(ctx context.Context, id int) (*entities.Municipality, error) {
	query := `
		SELECT id, name, state, active, created_at, updated_at
		FROM municipalities
		WHERE id = $1
	`
	
	municipality := &entities.Municipality{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&municipality.ID,
		&municipality.Name,
		&municipality.State,
		&municipality.Active,
		&municipality.CreatedAt,
		&municipality.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error getting municipality: %w", err)
	}
	
	return municipality, nil
}

func (r *municipalityRepository) GetByName(ctx context.Context, name string) (*entities.Municipality, error) {
	query := `
		SELECT id, name, state, active, created_at, updated_at
		FROM municipalities
		WHERE name = $1
	`
	
	municipality := &entities.Municipality{}
	err := r.db.Pool.QueryRow(ctx, query, name).Scan(
		&municipality.ID,
		&municipality.Name,
		&municipality.State,
		&municipality.Active,
		&municipality.CreatedAt,
		&municipality.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error getting municipality: %w", err)
	}
	
	return municipality, nil
}

func (r *municipalityRepository) List(ctx context.Context) ([]*entities.Municipality, error) {
	query := `
		SELECT id, name, state, active, created_at, updated_at
		FROM municipalities
		WHERE active = true
		ORDER BY name ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing municipalities: %w", err)
	}
	defer rows.Close()
	
	var municipalities []*entities.Municipality
	for rows.Next() {
		municipality := &entities.Municipality{}
		err := rows.Scan(
			&municipality.ID,
			&municipality.Name,
			&municipality.State,
			&municipality.Active,
			&municipality.CreatedAt,
			&municipality.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning municipality: %w", err)
		}
		municipalities = append(municipalities, municipality)
	}
	
	return municipalities, nil
}

func (r *municipalityRepository) Update(ctx context.Context, municipality *entities.Municipality) error {
	query := `
		UPDATE municipalities
		SET name = $2, state = $3, active = $4, updated_at = NOW()
		WHERE id = $1
	`
	
	cmdTag, err := r.db.Pool.Exec(ctx, query, 
		municipality.ID,
		municipality.Name,
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

func (r *municipalityRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM municipalities WHERE id = $1`
	
	cmdTag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting municipality: %w", err)
	}
	
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("municipality not found")
	}
	
	return nil
}