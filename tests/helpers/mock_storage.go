package helpers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
	"github.com/luizhreis/domain-watcher/internal/storage"
)

type MockStorage struct {
	// Add fields as needed for your mock storage implementation
	domains                 map[uuid.UUID]*models.Domain
	callHistory             []string
	createDomainShouldError bool
}

var _ storage.Storage = (*MockStorage)(nil)

func NewMockStorage() *MockStorage {
	return &MockStorage{
		domains:     make(map[uuid.UUID]*models.Domain),
		callHistory: []string{},
	}
}

func (m *MockStorage) CreateDomain(domain *models.Domain) error {
	m.callHistory = append(m.callHistory, "CreateDomain")

	if m.createDomainShouldError {
		return errors.New("simulated CreateDomain error")
	}

	if domain.ID == uuid.Nil {
		domain.ID = uuid.New()
	}
	m.domains[domain.ID] = domain
	return nil
}

func (m *MockStorage) GetDomain(id uuid.UUID) (*models.Domain, error) {
	m.callHistory = append(m.callHistory, "GetDomain")
	domain, exists := m.domains[id]
	if !exists {
		return nil, storage.ErrDomainNotFound
	}
	return domain, nil
}

func (m *MockStorage) GetAllDomains() ([]*models.Domain, error) {
	m.callHistory = append(m.callHistory, "GetAllDomains")
	domains := make([]*models.Domain, 0, len(m.domains))
	for _, domain := range m.domains {
		domains = append(domains, domain)
	}
	return domains, nil
}

func (m *MockStorage) UpdateDomain(domain *models.Domain) error {
	m.callHistory = append(m.callHistory, "UpdateDomain")
	if _, exists := m.domains[domain.ID]; !exists {
		return storage.ErrDomainNotFound
	}
	m.domains[domain.ID] = domain
	return nil
}

func (m *MockStorage) DeleteDomain(id uuid.UUID) error {
	m.callHistory = append(m.callHistory, "DeleteDomain")
	if _, exists := m.domains[id]; !exists {
		return storage.ErrDomainNotFound
	}
	delete(m.domains, id)
	return nil
}

func (m *MockStorage) GetCallHistory() []string {
	return m.callHistory
}

func (m *MockStorage) ClearCallHistory() {
	m.callHistory = []string{}
}

func (m *MockStorage) Reset() {
	m.domains = make(map[uuid.UUID]*models.Domain)
	m.createDomainShouldError = false
	m.ClearCallHistory()
}

func (m *MockStorage) SetCreateDomainError(shouldError bool) {
	m.createDomainShouldError = shouldError
}
