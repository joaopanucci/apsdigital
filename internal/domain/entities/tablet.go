package entities

import (
	"time"

	"github.com/google/uuid"
)

type TabletStatus string

const (
	TabletStatusActive      TabletStatus = "ativo"
	TabletStatusReturned    TabletStatus = "devolvido"
	TabletStatusBroken      TabletStatus = "quebrado"
	TabletStatusStolen      TabletStatus = "furtado"
	TabletStatusMaintenance TabletStatus = "manutencao"
	TabletStatusAvailable   TabletStatus = "disponivel"
	TabletStatusAssigned    TabletStatus = "atribuido"
)

type Tablet struct {
	ID             int          `json:"id" db:"id"`
	AssignedTo     *uuid.UUID   `json:"assigned_to" db:"assigned_to"`
	AssignedUserID *uuid.UUID   `json:"assigned_user_id" db:"assigned_user_id"`
	UserCPF        *string      `json:"user_cpf" db:"user_cpf"`
	MunicipalityID int          `json:"municipality_id" db:"municipality_id"`
	Status         TabletStatus `json:"status" db:"status"`
	AssetCode      string       `json:"asset_code" db:"asset_code"`
	Model          string       `json:"model" db:"model"`
	SerialNumber   string       `json:"serial_number" db:"serial_number"`
	AssignedAt     *time.Time   `json:"assigned_at" db:"assigned_at"`
	ReturnedAt     *time.Time   `json:"returned_at" db:"returned_at"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`

	// Relations
	AssignedUser *User         `json:"assigned_user,omitempty"`
	Municipality *Municipality `json:"municipality,omitempty"`
}
