package domain

import (
	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
)

type Domain interface {
	Create(domain *models.Domain) (uuid.UUID, error)
	Get(id uuid.UUID) (*models.Domain, error)
	List(page, pageSize int) ([]*models.Domain, error)
	Update(domain *models.Domain) error
	Delete(id uuid.UUID) error
}
