package models

import (
	"time"

	"github.com/google/uuid"
)

type CheckResult struct {
	ID            uuid.UUID `json:"id" db:"id"`
	DomainID      uuid.UUID `json:"domain_id" db:"domain_id"`
	StatusCode    int       `json:"status_code" db:"status_code"`
	ResponseTime  int64     `json:"response_time_ms" db:"response_time_ms"`
	Error         string    `json:"error,omitempty" db:"error"`
	RedirectURL   string    `json:"redirect_url,omitempty" db:"redirect_url"`
	RedirectCount int       `json:"redirect_count" db:"redirect_count"`
	CheckedAt     time.Time `json:"checked_at" db:"checked_at"`
	ContentLength int64     `json:"content_length" db:"content_length"`
	Server        string    `json:"server,omitempty" db:"server"`
	ResolvedIP    string    `json:"resolved_ip,omitempty" db:"resolved_ip"`
}
