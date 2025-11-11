package repositories

import (
	"context"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/infra/db"
)

type ProfessionRepository interface {
	Create(ctx context.Context, profession *entities.Profession) error
	GetByID(ctx context.Context, id int) (*entities.Profession, error)
	GetByName(ctx context.Context, name string) (*entities.Profession, error)
	GetAll(ctx context.Context) ([]*entities.Profession, error)
	Update(ctx context.Context, profession *entities.Profession) error
	Delete(ctx context.Context, id int) error
}

type professionRepository struct {
	db *db.PostgresDB
}

func NewProfessionRepository(db *db.PostgresDB) ProfessionRepository {
	return &professionRepository{
		db: db,
	}
}

func (r *professionRepository) Create(ctx context.Context, profession *entities.Profession) error {
	query := `
		INSERT INTO professions (name, created_at, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	
	row := r.db.Pool.QueryRow(ctx, query, profession.Name)
	return row.Scan(&profession.ID, &profession.CreatedAt, &profession.UpdatedAt)
}

func (r *professionRepository) GetByID(ctx context.Context, id int) (*entities.Profession, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM professions
		WHERE id = $1
	`
	
	var profession entities.Profession
	row := r.db.Pool.QueryRow(ctx, query, id)
	
	err := row.Scan(
		&profession.ID,
		&profession.Name,
		&profession.CreatedAt,
		&profession.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &profession, nil
}

func (r *professionRepository) GetByName(ctx context.Context, name string) (*entities.Profession, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM professions
		WHERE name = $1
	`
	
	var profession entities.Profession
	row := r.db.Pool.QueryRow(ctx, query, name)
	
	err := row.Scan(
		&profession.ID,
		&profession.Name,
		&profession.CreatedAt,
		&profession.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &profession, nil
}

func (r *professionRepository) GetAll(ctx context.Context) ([]*entities.Profession, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM professions
		ORDER BY name ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var professions []*entities.Profession
	for rows.Next() {
		var profession entities.Profession
		err := rows.Scan(
			&profession.ID,
			&profession.Name,
			&profession.CreatedAt,
			&profession.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		professions = append(professions, &profession)
	}
	
	return professions, rows.Err()
}

func (r *professionRepository) Update(ctx context.Context, profession *entities.Profession) error {
	query := `
		UPDATE professions
		SET name = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	row := r.db.Pool.QueryRow(ctx, query, profession.ID, profession.Name)
	return row.Scan(&profession.UpdatedAt)
}

func (r *professionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM professions WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}