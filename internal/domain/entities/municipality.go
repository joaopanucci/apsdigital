package entities

import (
	"time"
)

type Municipality struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	IBGECode  string    `json:"ibge_code" db:"ibge_code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}