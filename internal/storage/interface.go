package storage

import (
	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
)

type Storage interface {
	CreateDomain(domain *models.Domain) (uuid.UUID, error)
	GetDomain(id uuid.UUID) (*models.Domain, error)
	GetAllDomains() ([]*models.Domain, error)
	UpdateDomain(domain *models.Domain) error
	DeleteDomain(id uuid.UUID) error
}
