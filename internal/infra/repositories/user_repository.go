package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/infra/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *db.PostgresDB
}

func NewUserRepository(db *db.PostgresDB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, password, name, cpf, phone, role_id, profession_id, municipality, unit, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	user.ID = uuid.New()
	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.Name, user.CPF,
		user.Phone, user.RoleID, user.ProfessionID, user.Municipality,
		user.Unit, user.Status,
	)
	
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.password, u.name, u.cpf, u.phone, 
		       u.role_id, u.profession_id, u.municipality, u.unit, u.status, 
		       u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.level, r.created_at, r.updated_at,
		       p.id, p.name, p.created_at, p.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN professions p ON u.profession_id = p.id
		WHERE u.id = $1
	`
	
	row := r.db.Pool.QueryRow(ctx, query, id)
	
	var user entities.User
	var role entities.Role
	var profession entities.Profession
	
	err := row.Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.CPF,
		&user.Phone, &user.RoleID, &user.ProfessionID, &user.Municipality,
		&user.Unit, &user.Status, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt, &role.UpdatedAt,
		&profession.ID, &profession.Name, &profession.Description, &profession.IsActive, &profession.CreatedAt, &profession.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	user.Role = &role
	user.Profession = &profession
	
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.password, u.name, u.cpf, u.phone, 
		       u.role_id, u.profession_id, u.municipality, u.unit, u.status, 
		       u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.level, r.created_at, r.updated_at,
		       p.id, p.name, p.created_at, p.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN professions p ON u.profession_id = p.id
		WHERE u.email = $1
	`
	
	row := r.db.Pool.QueryRow(ctx, query, email)
	
	var user entities.User
	var role entities.Role
	var profession entities.Profession
	
	err := row.Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.CPF,
		&user.Phone, &user.RoleID, &user.ProfessionID, &user.Municipality,
		&user.Unit, &user.Status, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt, &role.UpdatedAt,
		&profession.ID, &profession.Name, &profession.Description, &profession.IsActive, &profession.CreatedAt, &profession.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	user.Role = &role
	user.Profession = &profession
	
	return &user, nil
}

func (r *userRepository) GetByCPF(ctx context.Context, cpf string) (*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.password, u.name, u.cpf, u.phone, 
		       u.role_id, u.profession_id, u.municipality, u.unit, u.status, 
		       u.created_at, u.updated_at
		FROM users u
		WHERE u.cpf = $1
	`
	
	row := r.db.Pool.QueryRow(ctx, query, cpf)
	
	var user entities.User
	
	err := row.Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.CPF,
		&user.Phone, &user.RoleID, &user.ProfessionID, &user.Municipality,
		&user.Unit, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $2, name = $3, phone = $4, role_id = $5, profession_id = $6, 
		    municipality = $7, unit = $8, status = $9, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.Phone, user.RoleID,
		user.ProfessionID, user.Municipality, user.Unit, user.Status,
	)
	
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.cpf, u.phone, 
		       u.role_id, u.profession_id, u.municipality, u.unit, u.status, 
		       u.created_at, u.updated_at,
		       r.name, p.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN professions p ON u.profession_id = p.id
	`
	
	var conditions []string
	var args []interface{}
	argCount := 0
	
	for key, value := range filters {
		argCount++
		switch key {
		case "municipality":
			conditions = append(conditions, fmt.Sprintf("u.municipality = $%d", argCount))
		case "status":
			conditions = append(conditions, fmt.Sprintf("u.status = $%d", argCount))
		case "role_id":
			conditions = append(conditions, fmt.Sprintf("u.role_id = $%d", argCount))
		case "profession_id":
			conditions = append(conditions, fmt.Sprintf("u.profession_id = $%d", argCount))
		}
		args = append(args, value)
	}
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += " ORDER BY u.created_at DESC"
	
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*entities.User
	
	for rows.Next() {
		var user entities.User
		var roleName, professionName string
		
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.CPF, &user.Phone,
			&user.RoleID, &user.ProfessionID, &user.Municipality, &user.Unit,
			&user.Status, &user.CreatedAt, &user.UpdatedAt,
			&roleName, &professionName,
		)
		if err != nil {
			return nil, err
		}
		
		user.Role = &entities.Role{Name: roleName}
		user.Profession = &entities.Profession{Name: professionName}
		
		users = append(users, &user)
	}
	
	return users, nil
}

func (r *userRepository) GetPendingAuthorization(ctx context.Context) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.cpf, u.phone, 
		       u.role_id, u.profession_id, u.municipality, u.unit, u.status, 
		       u.created_at, u.updated_at,
		       r.name, r.level, p.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN professions p ON u.profession_id = p.id
		WHERE u.status = 'pending_authorization'
		ORDER BY u.created_at ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*entities.User
	
	for rows.Next() {
		var user entities.User
		var roleName string
		var roleLevel int
		var professionName string
		
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.CPF, &user.Phone,
			&user.RoleID, &user.ProfessionID, &user.Municipality, &user.Unit,
			&user.Status, &user.CreatedAt, &user.UpdatedAt,
			&roleName, &roleLevel, &professionName,
		)
		if err != nil {
			return nil, err
		}
		
		user.Role = &entities.Role{Name: roleName, Level: roleLevel}
		user.Profession = &entities.Profession{Name: professionName}
		
		users = append(users, &user)
	}
	
	return users, nil
}
