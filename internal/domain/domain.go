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

func (d *domain) Get(id uuid.UUID) (*models.Domain, error) {
	if !isValidUUID(id) {
		return nil, ErrInvalidUUID
	}

	domain, err := d.storage.GetDomain(id)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (d *domain) List(page, pageSize int) ([]*models.Domain, error) {
	if page < 1 || pageSize < 1 {
		return nil, ErrInvalidPagination
	}

	domains, err := d.storage.ListDomains(page, pageSize)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func isValidUUID(id uuid.UUID) bool {
	return id != uuid.Nil
}
