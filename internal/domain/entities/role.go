package entities

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Level       int       `json:"level" db:"level"` // 1=ADM, 2=Coordenador, 3=Gerente, 4=ACS
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

const (
	RoleAdmin       = "ADM"
	RoleCoordenador = "Coordenador"
	RoleGerente     = "Gerente"
	RoleACS         = "ACS"
)

const (
	LevelAdmin       = 1
	LevelCoordenador = 2
	LevelGerente     = 3
	LevelACS         = 4
)