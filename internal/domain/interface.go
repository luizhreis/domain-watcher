package domain

import "github.com/luizhreis/domain-watcher/internal/models"

type Domain interface {
	Create(domain *models.Domain) error
}
