package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
)

type ResolutionRepository struct {
	db *sql.DB
}

func NewResolutionRepository(db *sql.DB) *ResolutionRepository {
	return &ResolutionRepository{db: db}
}

func (r *ResolutionRepository) Create(ctx context.Context, resolution *entities.Resolution) error {
	query := `
		INSERT INTO resolutions (title, file_url, competence, type, year, number, uploaded_by, municipality_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	row := r.db.QueryRowContext(ctx, query, 
		resolution.Title, resolution.FileURL, resolution.Competence, 
		resolution.Type, resolution.Year, resolution.Number, 
		resolution.UploadedBy, resolution.MunicipalityID)
	
	return row.Scan(&resolution.ID, &resolution.CreatedAt, &resolution.UpdatedAt)
}

func (r *ResolutionRepository) GetByID(ctx context.Context, id uint) (*entities.Resolution, error) {
	query := `
		SELECT r.id, r.title, r.file_url, r.competence, r.type, r.year, r.number, r.uploaded_by, r.municipality_id, r.created_at, r.updated_at,
		       u.name as uploaded_by_name, u.cpf as uploaded_by_cpf,
		       m.name as municipality_name
		FROM resolutions r
		LEFT JOIN users u ON r.uploaded_by = u.id
		LEFT JOIN municipalities m ON r.municipality_id = m.id
		WHERE r.id = $1
	`
	
	var resolution entities.Resolution
	row := r.db.QueryRowContext(ctx, query, id)
	
	err := row.Scan(
		&resolution.ID, &resolution.Title, &resolution.FileURL, &resolution.Competence, &resolution.Type,
		&resolution.Year, &resolution.Number, &resolution.UploadedBy, &resolution.MunicipalityID,
		&resolution.CreatedAt, &resolution.UpdatedAt, &resolution.UploadedByName, &resolution.UploadedByCPF, &resolution.MunicipalityName,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &resolution, nil
}

func (r *ResolutionRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Resolution, error) {
	query := `
		SELECT r.id, r.title, r.file_url, r.competence, r.type, r.year, r.number, r.uploaded_by, r.municipality_id, r.created_at, r.updated_at,
		       u.name as uploaded_by_name, u.cpf as uploaded_by_cpf,
		       m.name as municipality_name
		FROM resolutions r
		LEFT JOIN users u ON r.uploaded_by = u.id
		LEFT JOIN municipalities m ON r.municipality_id = m.id
	`
	
	var whereConditions []string
	var args []interface{}
	argIndex := 1
	
	// Filter by municipality
	if municipalityID, ok := filters["municipality_id"]; ok && municipalityID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("r.municipality_id = $%d", argIndex))
		args = append(args, municipalityID)
		argIndex++
	}
	
	// Filter by year
	if year, ok := filters["year"]; ok && year != nil && year != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("r.year = $%d", argIndex))
		args = append(args, year)
		argIndex++
	}
	
	// Filter by type
	if resType, ok := filters["type"]; ok && resType != nil && resType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("r.type = $%d", argIndex))
		args = append(args, resType)
		argIndex++
	}
	
	// Filter by number
	if number, ok := filters["number"]; ok && number != nil && number != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("r.number ILIKE $%d", argIndex))
		args = append(args, "%"+fmt.Sprintf("%v", number)+"%")
		argIndex++
	}
	
	// Filter by competence
	if competence, ok := filters["competence"]; ok && competence != nil && competence != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("r.competence = $%d", argIndex))
		args = append(args, competence)
		argIndex++
	}
	
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}
	
	query += " ORDER BY r.created_at DESC"
	
	// Add limit if specified
	if limit, ok := filters["limit"]; ok && limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}
	
	// Add offset if specified
	if offset, ok := filters["offset"]; ok && offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var resolutions []*entities.Resolution
	for rows.Next() {
		var resolution entities.Resolution
		err := rows.Scan(
			&resolution.ID, &resolution.Title, &resolution.FileURL, &resolution.Competence, &resolution.Type,
			&resolution.Year, &resolution.Number, &resolution.UploadedBy, &resolution.MunicipalityID,
			&resolution.CreatedAt, &resolution.UpdatedAt, &resolution.UploadedByName, &resolution.UploadedByCPF, &resolution.MunicipalityName,
		)
		if err != nil {
			return nil, err
		}
		resolutions = append(resolutions, &resolution)
	}
	
	return resolutions, nil
}

func (r *ResolutionRepository) Update(ctx context.Context, resolution *entities.Resolution) error {
	query := `
		UPDATE resolutions 
		SET title = $2, file_url = $3, competence = $4, type = $5, year = $6, number = $7, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		resolution.ID, resolution.Title, resolution.FileURL, resolution.Competence, 
		resolution.Type, resolution.Year, resolution.Number)
	return err
}

func (r *ResolutionRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM resolutions WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ResolutionRepository) GetTypes(ctx context.Context, municipalityID *uint) ([]string, error) {
	query := `
		SELECT DISTINCT type 
		FROM resolutions 
	`
	
	var args []interface{}
	if municipalityID != nil {
		query += " WHERE municipality_id = $1"
		args = append(args, *municipalityID)
	}
	
	query += " ORDER BY type ASC"
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var types []string
	for rows.Next() {
		var resType string
		if err := rows.Scan(&resType); err != nil {
			return nil, err
		}
		types = append(types, resType)
	}
	
	return types, nil
}

func (r *ResolutionRepository) GetYears(ctx context.Context, municipalityID *uint) ([]int, error) {
	query := `
		SELECT DISTINCT year 
		FROM resolutions 
	`
	
	var args []interface{}
	if municipalityID != nil {
		query += " WHERE municipality_id = $1"
		args = append(args, *municipalityID)
	}
	
	query += " ORDER BY year DESC"
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var years []int
	for rows.Next() {
		var year int
		if err := rows.Scan(&year); err != nil {
			return nil, err
		}
		years = append(years, year)
	}
	
	return years, nil
}

func (r *ResolutionRepository) GetRecent(ctx context.Context, municipalityID *uint, limit int) ([]*entities.Resolution, error) {
	query := `
		SELECT r.id, r.title, r.file_url, r.competence, r.type, r.year, r.number, r.uploaded_by, r.municipality_id, r.created_at, r.updated_at,
		       u.name as uploaded_by_name, u.cpf as uploaded_by_cpf,
		       m.name as municipality_name
		FROM resolutions r
		LEFT JOIN users u ON r.uploaded_by = u.id
		LEFT JOIN municipalities m ON r.municipality_id = m.id
	`
	
	var args []interface{}
	if municipalityID != nil {
		query += " WHERE r.municipality_id = $1"
		args = append(args, *municipalityID)
		query += fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d", len(args)+1)
		args = append(args, limit)
	} else {
		query += " ORDER BY r.created_at DESC LIMIT $1"
		args = append(args, limit)
	}
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var resolutions []*entities.Resolution
	for rows.Next() {
		var resolution entities.Resolution
		err := rows.Scan(
			&resolution.ID, &resolution.Title, &resolution.FileURL, &resolution.Competence, &resolution.Type,
			&resolution.Year, &resolution.Number, &resolution.UploadedBy, &resolution.MunicipalityID,
			&resolution.CreatedAt, &resolution.UpdatedAt, &resolution.UploadedByName, &resolution.UploadedByCPF, &resolution.MunicipalityName,
		)
		if err != nil {
			return nil, err
		}
		resolutions = append(resolutions, &resolution)
	}
	
	return resolutions, nil
}