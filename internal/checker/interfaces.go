package checker

import "github.com/luizhreis/domain-watcher/internal/models"

type Checker interface {
	CheckDomain(domain *models.Domain) (*models.CheckResult, error)
}
