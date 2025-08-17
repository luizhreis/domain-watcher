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
	getDomainShouldError    bool
	listDomainsShouldError  bool
}

var _ storage.Storage = (*MockStorage)(nil)

func NewMockStorage() *MockStorage {
	return &MockStorage{
		domains:     make(map[uuid.UUID]*models.Domain),
		callHistory: []string{},
	}
}

func (m *MockStorage) CreateDomain(domain *models.Domain) (uuid.UUID, error) {
	m.callHistory = append(m.callHistory, "CreateDomain")

	if m.createDomainShouldError {
		return uuid.Nil, errors.New("simulated CreateDomain error")
	}

	if domain.ID == uuid.Nil {
		domain.ID = uuid.New()
	}
	m.domains[domain.ID] = domain
	return domain.ID, nil
}

func (m *MockStorage) GetDomain(id uuid.UUID) (*models.Domain, error) {
	m.callHistory = append(m.callHistory, "GetDomain")

	if m.getDomainShouldError {
		return nil, errors.New("mock get domain error")
	}

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

func (m *MockStorage) ListDomains(page, pageSize int) ([]*models.Domain, error) {
	m.callHistory = append(m.callHistory, "ListDomains")

	if m.listDomainsShouldError {
		return nil, errors.New("mock list domains error")
	}

	if page < 1 || pageSize < 1 {
		return nil, errors.New("page and pageSize must be greater than 0")
	}

	// Converte map para slice para permitir paginação
	allDomains := make([]*models.Domain, 0, len(m.domains))
	for _, domain := range m.domains {
		allDomains = append(allDomains, domain)
	}

	// Calcula offset
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Verifica limites
	if startIndex >= len(allDomains) {
		return []*models.Domain{}, nil // Página vazia, mas não é erro
	}

	if endIndex > len(allDomains) {
		endIndex = len(allDomains)
	}

	return allDomains[startIndex:endIndex], nil
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
	m.getDomainShouldError = false
	m.listDomainsShouldError = false
	m.ClearCallHistory()
}

func (m *MockStorage) SetCreateDomainError(shouldError bool) {
	m.createDomainShouldError = shouldError
}

func (m *MockStorage) SetGetDomainError(shouldError bool) {
	m.getDomainShouldError = shouldError
}

func (m *MockStorage) SetListDomainsError(shouldError bool) {
	m.listDomainsShouldError = shouldError
}
