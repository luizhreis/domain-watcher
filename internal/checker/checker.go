package checker

import (
	"time"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/dns"
	"github.com/luizhreis/domain-watcher/internal/models"
)

type checker struct {
	DNSResolver dns.DNS
}

func NewChecker(dnsResolver dns.DNS) Checker {
	return &checker{
		DNSResolver: dnsResolver,
	}
}

var _ Checker = (*checker)(nil)

func (c *checker) CheckDomain(domain *models.Domain) (*models.CheckResult, error) {
	timestamp := time.Now()

	resolvedIP, err := c.DNSResolver.Resolve(domain.URL)
	if err != nil {
		return nil, err
	}

	return &models.CheckResult{
		ID:            uuid.New(),
		DomainID:      domain.ID,
		StatusCode:    200,
		ResponseTime:  100,
		Error:         "",
		RedirectURL:   "",
		RedirectCount: 0,
		CheckedAt:     timestamp,
		ContentLength: 1024,
		Server:        "nginx",
		ResolvedIP:    resolvedIP,
	}, nil
}
