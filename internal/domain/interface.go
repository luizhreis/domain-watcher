package domain

import "github.com/luizhreis/domain-watcher/internal/models"

type Domain interface {
	CreateDomain(domain *models.Domain) error
}
