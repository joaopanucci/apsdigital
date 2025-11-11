package entities

import (
	"time"

	"github.com/google/uuid"
)

type Profession struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Predefined professions
const (
	ProfessionEnfermeiro              = "Enfermeiro"
	ProfessionMedico                  = "Médico"
	ProfessionAgenteComunitario       = "Agente Comunitário de Saúde"
	ProfessionOdontologo             = "Odontólogo"
	ProfessionTecnicoEnfermagem      = "Técnico de Enfermagem"
	ProfessionFisioterapeuta         = "Fisioterapeuta"
	ProfessionPsicologo              = "Psicólogo"
	ProfessionNutricionista          = "Nutricionista"
	ProfessionFarmaceutico           = "Farmacêutico"
	ProfessionAssistenteSocial       = "Assistente Social"
)