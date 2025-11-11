package entities

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	FileURL          string     `json:"file_url" db:"file_url"`
	CompetenceMonth  string     `json:"competence_month" db:"competence_month"` // YYYY-MM format
	CompetenceDate   *time.Time `json:"competence_date" db:"competence_date"`
	CompetenceYear   int        `json:"competence_year" db:"competence_year"`
	Municipality     string     `json:"municipality" db:"municipality"` // Keep for backward compatibility
	MunicipalityID   *int       `json:"municipality_id" db:"municipality_id"`
	UploadedBy       uuid.UUID  `json:"uploaded_by" db:"uploaded_by"`
	OriginalFileName string     `json:"original_file_name" db:"original_file_name"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	UploadedByUser *User         `json:"uploaded_by_user,omitempty"`
	MunicipalityInfo *Municipality `json:"municipality_info,omitempty"`
}