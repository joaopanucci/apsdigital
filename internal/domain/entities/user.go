package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive               UserStatus = "active"
	UserStatusPendingAuthorization UserStatus = "pending_authorization"
	UserStatusBlocked              UserStatus = "blocked"
	UserStatusInactive             UserStatus = "inactive"
)

type User struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Email          string     `json:"email" db:"email"`
	Password       string     `json:"-" db:"password"`
	Name           string     `json:"name" db:"name"`
	CPF            string     `json:"cpf" db:"cpf"`
	Phone          string     `json:"phone" db:"phone"`
	RoleID         uuid.UUID  `json:"role_id" db:"role_id"`
	ProfessionID   *int       `json:"profession_id" db:"profession_id"`
	Municipality   string     `json:"municipality" db:"municipality"` // Keep for backward compatibility
	MunicipalityID *int       `json:"municipality_id" db:"municipality_id"`
	Unit           string     `json:"unit" db:"unit"`
	Status         UserStatus `json:"status" db:"status"`
	IsAuthorized   bool       `json:"is_authorized" db:"is_authorized"`
	ProfilePhoto   string     `json:"profile_photo" db:"profile_photo"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`

	// Relations
	Role             *Role         `json:"role,omitempty"`
	Profession       *Profession   `json:"profession,omitempty"`
	MunicipalityInfo *Municipality `json:"municipality_info,omitempty"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked"`
}
