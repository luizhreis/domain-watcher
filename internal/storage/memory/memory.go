package memory

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
)

var (
	// ErrDomainNotFound é retornado quando um domínio não é encontrado
	ErrDomainNotFound = errors.New("domain not found")
)

// MemoryStorage é uma implementação in-memory do Storage
type MemoryStorage struct {
	domains map[uuid.UUID]*models.Domain
}

// NewMemoryStorage cria uma nova instância de MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		domains: make(map[uuid.UUID]*models.Domain),
	}
}

// CreateDomain cria um novo domínio no storage
func (m *MemoryStorage) CreateDomain(domain *models.Domain) (uuid.UUID, error) {
	domain.ID = uuid.New()

	// Define timestamps
	now := time.Now()
	domain.CreatedAt = now
	domain.UpdatedAt = now

	// Armazena no map
	m.domains[domain.ID] = domain

	return domain.ID, nil
}

// GetDomain busca um domínio pelo ID
func (m *MemoryStorage) GetDomain(id uuid.UUID) (*models.Domain, error) {
	domain, exists := m.domains[id]
	if !exists {
		return nil, ErrDomainNotFound
	}
	return domain, nil
}

// GetAllDomains retorna todos os domínios
func (m *MemoryStorage) GetAllDomains() ([]*models.Domain, error) {
	domains := make([]*models.Domain, 0, len(m.domains))
	for _, domain := range m.domains {
		domains = append(domains, domain)
	}
	return domains, nil
}

// ListDomains retorna domínios com paginação
func (m *MemoryStorage) ListDomains(page, pageSize int) ([]*models.Domain, error) {
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

// UpdateDomain atualiza um domínio existente
func (m *MemoryStorage) UpdateDomain(domain *models.Domain) error {
	if _, exists := m.domains[domain.ID]; !exists {
		return ErrDomainNotFound
	}

	domain.UpdatedAt = time.Now()
	m.domains[domain.ID] = domain

	return nil
}

// DeleteDomain remove um domínio
func (m *MemoryStorage) DeleteDomain(id uuid.UUID) error {
	if _, exists := m.domains[id]; !exists {
		return ErrDomainNotFound
	}

	delete(m.domains, id)
	return nil
}
