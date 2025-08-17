package domain

import (
	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
)

type Domain interface {
	Create(domain *models.Domain) (uuid.UUID, error)
}
