package entities

import (
	"time"

	"github.com/google/uuid"
)

type TabletRequestType string

const (
	TabletRequestTypeNew        TabletRequestType = "novo"
	TabletRequestTypeReturn     TabletRequestType = "devolucao"
	TabletRequestTypeBreakage   TabletRequestType = "quebra"
	TabletRequestTypeTheft      TabletRequestType = "furto"
)

type TabletRequestStatus string

const (
	TabletRequestStatusPending    TabletRequestStatus = "pending"
	TabletRequestStatusApproved   TabletRequestStatus = "approved"
	TabletRequestStatusRejected   TabletRequestStatus = "rejected"
	TabletRequestStatusCompleted  TabletRequestStatus = "completed"
)

type TabletRequest struct {
	ID             uuid.UUID           `json:"id" db:"id"`
	UserID         uuid.UUID           `json:"user_id" db:"user_id"`
	Type           TabletRequestType   `json:"type" db:"type"`
	Status         TabletRequestStatus `json:"status" db:"status"`
	Justification  string              `json:"justification" db:"justification"`
	Description    string              `json:"description" db:"description"`
	Photos         []string            `json:"photos" db:"photos"` // JSON array of file paths
	DocumentURL    string              `json:"document_url" db:"document_url"` // For BO in theft cases
	ApprovedBy     *uuid.UUID          `json:"approved_by" db:"approved_by"`
	ApprovedAt     *time.Time          `json:"approved_at" db:"approved_at"`
	RejectedBy     *uuid.UUID          `json:"rejected_by" db:"rejected_by"`
	RejectedAt     *time.Time          `json:"rejected_at" db:"rejected_at"`
	RejectionReason string             `json:"rejection_reason" db:"rejection_reason"`
	CreatedAt      time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" db:"updated_at"`
	
	// Relations
	User       *User `json:"user,omitempty"`
	ApprovedByUser *User `json:"approved_by_user,omitempty"`
	RejectedByUser *User `json:"rejected_by_user,omitempty"`
}
