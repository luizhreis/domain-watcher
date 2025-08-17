package models

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	URL       string    `json:"url" db:"url"`
	Timeout   int       `json:"timeout" db:"timeout"`
	IP        string    `json:"ip,omitempty" db:"ip"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
