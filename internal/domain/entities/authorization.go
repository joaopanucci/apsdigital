package entities

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationStatus string

const (
	AuthorizationStatusPending  AuthorizationStatus = "pending"
	AuthorizationStatusApproved AuthorizationStatus = "approved"
	AuthorizationStatusRejected AuthorizationStatus = "rejected"
)

type Authorization struct {
	ID              uuid.UUID           `json:"id" db:"id"`
	UserID          uuid.UUID           `json:"user_id" db:"user_id"`
	AuthorizedBy    *uuid.UUID          `json:"authorized_by" db:"authorized_by"`
	Status          AuthorizationStatus `json:"status" db:"status"`
	Comments        string              `json:"comments" db:"comments"`
	AuthorizedAt    *time.Time          `json:"authorized_at" db:"authorized_at"`
	RejectedAt      *time.Time          `json:"rejected_at" db:"rejected_at"`
	RejectionReason string              `json:"rejection_reason" db:"rejection_reason"`
	CreatedAt       time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" db:"updated_at"`
	
	// Relations
	User           *User `json:"user,omitempty"`
	AuthorizedByUser *User `json:"authorized_by_user,omitempty"`
}
