package entities

import (
	"time"

	"github.com/google/uuid"
)

type ResolutionType string

const (
	ResolutionTypeMS  ResolutionType = "MS"  // Ministério da Saúde
	ResolutionTypeSES ResolutionType = "SES" // Secretaria Estadual de Saúde
)

type Resolution struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Title            string         `json:"title" db:"title"`
	FileURL          string         `json:"file_url" db:"file_url"`
	Year             int            `json:"year" db:"year"`
	Type             ResolutionType `json:"type" db:"type"`
	Competence       string         `json:"competence" db:"competence"`
	Description      string         `json:"description" db:"description"`
	OriginalFileName string         `json:"original_file_name" db:"original_file_name"`
	FileSize         int64          `json:"file_size" db:"file_size"`
	UploadedBy       uuid.UUID      `json:"uploaded_by" db:"uploaded_by"`
	MunicipalityID   *int           `json:"municipality_id" db:"municipality_id"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
	
	// Relations
	UploadedByUser   *User         `json:"uploaded_by_user,omitempty"`
	MunicipalityInfo *Municipality `json:"municipality_info,omitempty"`
}