package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/infra/db"

	"github.com/google/uuid"
)

type roleRepository struct {
	db *db.PostgresDB
}

func NewRoleRepository(db *db.PostgresDB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entities.Role) error {
	query := `
		INSERT INTO roles (id, name, description, level)
		VALUES ($1, $2, $3, $4)
	`

	role.ID = uuid.New()
	_, err := r.db.Pool.Exec(ctx, query, role.ID, role.Name, role.Description, role.Level)
	return err
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	query := `
		SELECT id, name, description, level, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	row := r.db.Pool.QueryRow(ctx, query, id)

	var role entities.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entities.Role, error) {
	query := `
		SELECT id, name, description, level, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	row := r.db.Pool.QueryRow(ctx, query, name)

	var role entities.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]*entities.Role, error) {
	query := `
		SELECT id, name, description, level, created_at, updated_at
		FROM roles
		ORDER BY level ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entities.Role

	for rows.Next() {
		var role entities.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *roleRepository) Update(ctx context.Context, role *entities.Role) error {
	query := `
		UPDATE roles 
		SET name = $2, description = $3, level = $4, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, role.ID, role.Name, role.Description, role.Level)
	return err
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
