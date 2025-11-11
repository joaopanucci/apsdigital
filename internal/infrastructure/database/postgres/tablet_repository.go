package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"apsdigital/internal/domain/entities"
	"apsdigital/internal/domain/repositories"
)

type TabletPostgresRepository struct {
	db *pgxpool.Pool
}

func NewTabletRepository(db *pgxpool.Pool) repositories.TabletRepository {
	return &TabletPostgresRepository{db: db}
}

func (r *TabletPostgresRepository) Create(ctx context.Context, tablet *entities.Tablet) error {
	query := `
		INSERT INTO tablets (serial_number, model, status, assigned_user_id, user_cpf, assigned_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`
	
	err := r.db.QueryRow(ctx, query,
		tablet.SerialNumber,
		tablet.Model,
		tablet.Status,
		tablet.AssignedUserID,
		tablet.UserCPF,
		tablet.AssignedAt,
	).Scan(&tablet.ID)
	
	if err != nil {
		return fmt.Errorf("error creating tablet: %w", err)
	}
	
	return nil
}

func (r *TabletPostgresRepository) GetByID(ctx context.Context, id int) (*entities.Tablet, error) {
	query := `
		SELECT id, serial_number, model, status, assigned_user_id, user_cpf, assigned_at, created_at, updated_at
		FROM tablets
		WHERE id = $1
	`
	
	tablet := &entities.Tablet{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&tablet.ID,
		&tablet.SerialNumber,
		&tablet.Model,
		&tablet.Status,
		&tablet.AssignedUserID,
		&tablet.UserCPF,
		&tablet.AssignedAt,
		&tablet.CreatedAt,
		&tablet.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tablet not found")
		}
		return nil, fmt.Errorf("error getting tablet: %w", err)
	}
	
	return tablet, nil
}

func (r *TabletPostgresRepository) GetByUserCPF(ctx context.Context, cpf string) ([]*entities.Tablet, error) {
	query := `
		SELECT id, serial_number, model, status, assigned_user_id, user_cpf, assigned_at, created_at, updated_at
		FROM tablets
		WHERE user_cpf = $1
		ORDER BY assigned_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, cpf)
	if err != nil {
		return nil, fmt.Errorf("error getting tablets by user CPF: %w", err)
	}
	defer rows.Close()
	
	var tablets []*entities.Tablet
	for rows.Next() {
		tablet := &entities.Tablet{}
		err := rows.Scan(
			&tablet.ID,
			&tablet.SerialNumber,
			&tablet.Model,
			&tablet.Status,
			&tablet.AssignedUserID,
			&tablet.UserCPF,
			&tablet.AssignedAt,
			&tablet.CreatedAt,
			&tablet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tablet: %w", err)
		}
		tablets = append(tablets, tablet)
	}
	
	return tablets, nil
}

func (r *TabletPostgresRepository) GetByAssignedUser(ctx context.Context, userID uuid.UUID) ([]*entities.Tablet, error) {
	query := `
		SELECT id, serial_number, model, status, assigned_user_id, user_cpf, assigned_at, created_at, updated_at
		FROM tablets
		WHERE assigned_user_id = $1
		ORDER BY assigned_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting tablets by assigned user: %w", err)
	}
	defer rows.Close()
	
	var tablets []*entities.Tablet
	for rows.Next() {
		tablet := &entities.Tablet{}
		err := rows.Scan(
			&tablet.ID,
			&tablet.SerialNumber,
			&tablet.Model,
			&tablet.Status,
			&tablet.AssignedUserID,
			&tablet.UserCPF,
			&tablet.AssignedAt,
			&tablet.CreatedAt,
			&tablet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tablet: %w", err)
		}
		tablets = append(tablets, tablet)
	}
	
	return tablets, nil
}

func (r *TabletPostgresRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.Tablet, error) {
	query := `
		SELECT id, serial_number, model, status, assigned_user_id, user_cpf, assigned_at, created_at, updated_at
		FROM tablets
	`
	
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	if status, ok := filters["status"].(string); ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}
	
	if model, ok := filters["model"].(string); ok && model != "" {
		conditions = append(conditions, fmt.Sprintf("model ILIKE $%d", argIndex))
		args = append(args, "%"+model+"%")
		argIndex++
	}
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += " ORDER BY created_at DESC"
	
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing tablets: %w", err)
	}
	defer rows.Close()
	
	var tablets []*entities.Tablet
	for rows.Next() {
		tablet := &entities.Tablet{}
		err := rows.Scan(
			&tablet.ID,
			&tablet.SerialNumber,
			&tablet.Model,
			&tablet.Status,
			&tablet.AssignedUserID,
			&tablet.UserCPF,
			&tablet.AssignedAt,
			&tablet.CreatedAt,
			&tablet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tablet: %w", err)
		}
		tablets = append(tablets, tablet)
	}
	
	return tablets, nil
}

func (r *TabletPostgresRepository) Update(ctx context.Context, tablet *entities.Tablet) error {
	query := `
		UPDATE tablets
		SET serial_number = $2, model = $3, status = $4, assigned_user_id = $5, 
		    user_cpf = $6, assigned_at = $7, updated_at = NOW()
		WHERE id = $1
	`
	
	cmdTag, err := r.db.Exec(ctx, query,
		tablet.ID,
		tablet.SerialNumber,
		tablet.Model,
		tablet.Status,
		tablet.AssignedUserID,
		tablet.UserCPF,
		tablet.AssignedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error updating tablet: %w", err)
	}
	
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tablet not found")
	}
	
	return nil
}

func (r *TabletPostgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tablets WHERE id = $1`
	
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting tablet: %w", err)
	}
	
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tablet not found")
	}
	
	return nil
}