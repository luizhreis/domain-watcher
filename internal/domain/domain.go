package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
	"github.com/luizhreis/domain-watcher/internal/storage"
)

type domain struct {
	storage storage.Storage
}

var _ Domain = (*domain)(nil)

func NewDomain(storage storage.Storage) Domain {
	return &domain{
		storage: storage,
	}
}

func (d *domain) Create(domain *models.Domain) (uuid.UUID, error) {
	timestamp := time.Now()
	domain.CreatedAt = timestamp
	domain.UpdatedAt = timestamp

	id, err := d.storage.CreateDomain(domain)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
