package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
	"github.com/luizhreis/domain-watcher/tests/helpers"
)

func TestNewDomain(t *testing.T) {
	d := NewDomain(nil)
	if d == nil {
		t.Error("Expected NewDomain to return a non-nil instance")
	}
}

// TestCreateDomain testa a criação de um domínio (white-box)
func TestCreateDomain(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	d := &models.Domain{
		Name:    "test.com",
		URL:     "test.com",
		Timeout: 30,
	}

	id, err := domain.Create(d)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if id == uuid.Nil {
		t.Error("Expected a valid UUID, got nil")
	}

	// Verifica se o domínio foi criado
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to CreateDomain, got %d", len(storage.GetCallHistory()))
	}
}

// TestCreateDomainError testa o erro ao criar um domínio (white-box)
func TestCreateDomainError(t *testing.T) {
	storage := helpers.NewMockStorage()
	storage.SetCreateDomainError(true)
	domain := NewDomain(storage)

	d := &models.Domain{
		Name:    "error.com",
		URL:     "error.com",
		Timeout: 30,
	}

	if _, err := domain.Create(d); err == nil {
		t.Error("Expected error when creating domain, got nil")
	}

	// Verifica se o erro foi registrado
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to CreateDomain, got %d", len(storage.GetCallHistory()))
	}
}

// TestGetDomain testa a obtenção de um domínio (white-box)
func TestGetDomain(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	// Primeiro, cria um domínio para garantir que existe
	d := &models.Domain{
		Name:    "get.com",
		URL:     "get.com",
		Timeout: 30,
	}
	id, _ := domain.Create(d)

	retrieved, err := domain.Get(id)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved.ID != id {
		t.Errorf("Expected ID %v, got %v", id, retrieved.ID)
	}

	// Verifica se o domínio foi buscado
	if len(storage.GetCallHistory()) != 2 { // 1 para Create, 1 para Get
		t.Errorf("Expected 2 calls to storage, got %d", len(storage.GetCallHistory()))
	}
}

// TestGetDomainInvalidUUID testa o erro ao obter um domínio com UUID inválido (white-box)
func TestGetDomainInvalidUUID(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	invalidID := uuid.Nil // UUID inválido

	if _, err := domain.Get(invalidID); err == nil {
		t.Error("Expected error when getting domain with invalid UUID, got nil")
	}

	// Verifica se o erro foi registrado
	if len(storage.GetCallHistory()) != 0 { // Nenhuma chamada deve ter sido feita
		t.Errorf("Expected 0 calls to storage, got %d", len(storage.GetCallHistory()))
	}
}

// TestListDomains testa a listagem de domínios com paginação (white-box)
func TestListDomains(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	// Cria alguns domínios para testar a paginação
	domains := []*models.Domain{
		{Name: "test1.com", URL: "test1.com", Timeout: 30},
		{Name: "test2.com", URL: "test2.com", Timeout: 30},
		{Name: "test3.com", URL: "test3.com", Timeout: 30},
		{Name: "test4.com", URL: "test4.com", Timeout: 30},
		{Name: "test5.com", URL: "test5.com", Timeout: 30},
	}

	// Cria os domínios
	for _, d := range domains {
		_, err := domain.Create(d)
		if err != nil {
			t.Errorf("Failed to create domain: %v", err)
		}
	}

	// Testa paginação - primeira página com 2 itens
	page1, err := domain.List(1, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(page1) != 2 {
		t.Errorf("Expected 2 domains in page 1, got %d", len(page1))
	}

	// Testa paginação - segunda página com 2 itens
	page2, err := domain.List(2, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(page2) != 2 {
		t.Errorf("Expected 2 domains in page 2, got %d", len(page2))
	}

	// Testa paginação - terceira página com 1 item restante
	page3, err := domain.List(3, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(page3) != 1 {
		t.Errorf("Expected 1 domain in page 3, got %d", len(page3))
	}

	// Verifica chamadas ao storage (5 Create + 3 ListDomains)
	expectedCalls := 5 + 3
	if len(storage.GetCallHistory()) != expectedCalls {
		t.Errorf("Expected %d calls to storage, got %d", expectedCalls, len(storage.GetCallHistory()))
	}
}

// TestListDomainsInvalidPagination testa parâmetros de paginação inválidos (white-box)
func TestListDomainsInvalidPagination(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	tests := []struct {
		name     string
		page     int
		pageSize int
	}{
		{"Zero page", 0, 10},
		{"Negative page", -1, 10},
		{"Zero pageSize", 1, 0},
		{"Negative pageSize", 1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.List(tt.page, tt.pageSize)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
		})
	}

	// Nenhuma chamada ao storage deve ter sido feita
	if len(storage.GetCallHistory()) != 0 {
		t.Errorf("Expected 0 calls to storage, got %d", len(storage.GetCallHistory()))
	}
}

// TestListDomainsEmptyResult testa paginação com resultado vazio (white-box)
func TestListDomainsEmptyResult(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	// Não cria nenhum domínio, então resultado deve ser vazio

	result, err := domain.List(1, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d domains", len(result))
	}

	// Verifica se foi chamado o storage
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to ListDomains, got %d", len(storage.GetCallHistory()))
	}
}

// TestListDomainsStorageError testa erro do storage ao listar domínios (white-box)
func TestListDomainsStorageError(t *testing.T) {
	storage := helpers.NewMockStorage()
	storage.SetListDomainsError(true)
	domain := NewDomain(storage)

	_, err := domain.List(1, 10)
	if err == nil {
		t.Error("Expected error when storage fails, got nil")
	}

	// Verifica se o storage foi chamado
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to ListDomains, got %d", len(storage.GetCallHistory()))
	}
}

// TestGetDomainStorageError testa erro do storage ao buscar domínio (white-box)
func TestGetDomainStorageError(t *testing.T) {
	storage := helpers.NewMockStorage()
	storage.SetGetDomainError(true) // Você precisa adicionar este método no MockStorage
	domain := NewDomain(storage)

	validID := uuid.New()

	if _, err := domain.Get(validID); err == nil {
		t.Error("Expected error when storage fails, got nil")
	}

	// Verifica se o storage foi chamado
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to GetDomain, got %d", len(storage.GetCallHistory()))
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		id       uuid.UUID
		expected bool
	}{
		{
			name:     "Nil UUID should be invalid",
			id:       uuid.Nil,
			expected: false,
		},
		{
			name:     "Valid UUID should be valid",
			id:       uuid.New(),
			expected: true,
		},
		{
			name:     "Predefined UUID should be valid",
			id:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidUUID(tt.id)
			if result != tt.expected {
				t.Errorf("isValidUUID(%v) = %v, expected %v", tt.id, result, tt.expected)
			}
		})
	}
}

// TestUpdateDomain testa a atualização de um domínio (white-box)
func TestUpdateDomain(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	// Primeiro, cria um domínio
	d := &models.Domain{
		Name:    "update.com",
		URL:     "update.com",
		Timeout: 30,
	}
	id, err := domain.Create(d)
	if err != nil {
		t.Fatalf("Failed to create domain: %v", err)
	}

	// Modifica os dados
	d.Name = "updated.com"
	d.URL = "https://updated.com"
	d.Timeout = 60

	// Executa a atualização
	err = domain.Update(d)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se UpdatedAt foi modificado
	if d.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set after update")
	}

	// Verifica as chamadas ao storage (1 Create + 1 Update)
	if len(storage.GetCallHistory()) != 2 {
		t.Errorf("Expected 2 calls to storage, got %d", len(storage.GetCallHistory()))
	}

	// Verifica se foi o domínio correto que foi atualizado
	retrieved, err := domain.Get(id)
	if err != nil {
		t.Errorf("Failed to retrieve updated domain: %v", err)
	}

	if retrieved.Name != "updated.com" {
		t.Errorf("Expected Name 'updated.com', got '%s'", retrieved.Name)
	}
}

// TestUpdateDomainWithNil testa atualização com domain nil (white-box)
func TestUpdateDomainWithNil(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	err := domain.Update(nil)
	if err == nil {
		t.Error("Expected error when updating nil domain, got nil")
	}

	// Nenhuma chamada ao storage deve ter sido feita
	if len(storage.GetCallHistory()) != 0 {
		t.Errorf("Expected 0 calls to storage, got %d", len(storage.GetCallHistory()))
	}
}

// TestUpdateDomainWithInvalidUUID testa atualização com UUID inválido (white-box)
func TestUpdateDomainWithInvalidUUID(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	d := &models.Domain{
		ID:      uuid.Nil, // UUID inválido
		Name:    "invalid.com",
		URL:     "invalid.com",
		Timeout: 30,
	}

	err := domain.Update(d)
	if err == nil {
		t.Error("Expected error when updating domain with invalid UUID, got nil")
	}

	// Nenhuma chamada ao storage deve ter sido feita
	if len(storage.GetCallHistory()) != 0 {
		t.Errorf("Expected 0 calls to storage, got %d", len(storage.GetCallHistory()))
	}
}

// TestUpdateDomainStorageError testa erro do storage ao atualizar (white-box)
func TestUpdateDomainStorageError(t *testing.T) {
	storage := helpers.NewMockStorage()
	storage.SetUpdateDomainError(true)
	domain := NewDomain(storage)

	d := &models.Domain{
		ID:      uuid.New(),
		Name:    "error.com",
		URL:     "error.com",
		Timeout: 30,
	}

	err := domain.Update(d)
	if err == nil {
		t.Error("Expected error when storage fails, got nil")
	}

	// Verifica se o storage foi chamado
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to UpdateDomain, got %d", len(storage.GetCallHistory()))
	}
}

// TestDeleteDomain testa a exclusão de um domínio (white-box)
func TestDeleteDomain(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	// Primeiro, cria um domínio
	d := &models.Domain{
		Name:    "delete.com",
		URL:     "delete.com",
		Timeout: 30,
	}
	id, err := domain.Create(d)
	if err != nil {
		t.Fatalf("Failed to create domain: %v", err)
	}

	// Executa a exclusão
	err = domain.Delete(id)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica as chamadas ao storage (1 Create + 1 Delete)
	if len(storage.GetCallHistory()) != 2 {
		t.Errorf("Expected 2 calls to storage, got %d", len(storage.GetCallHistory()))
	}

	// Verifica se o domínio foi realmente removido
	_, err = domain.Get(id)
	if err == nil {
		t.Error("Expected error when getting deleted domain, got nil")
	}
}

// TestDeleteDomainWithInvalidUUID testa exclusão com UUID inválido (white-box)
func TestDeleteDomainWithInvalidUUID(t *testing.T) {
	storage := helpers.NewMockStorage()
	domain := NewDomain(storage)

	err := domain.Delete(uuid.Nil)
	if err == nil {
		t.Error("Expected error when deleting domain with invalid UUID, got nil")
	}

	// Nenhuma chamada ao storage deve ter sido feita
	if len(storage.GetCallHistory()) != 0 {
		t.Errorf("Expected 0 calls to storage, got %d", len(storage.GetCallHistory()))
	}
}

// TestDeleteDomainStorageError testa erro do storage ao excluir (white-box)
func TestDeleteDomainStorageError(t *testing.T) {
	storage := helpers.NewMockStorage()
	storage.SetDeleteDomainError(true)
	domain := NewDomain(storage)

	validID := uuid.New()

	err := domain.Delete(validID)
	if err == nil {
		t.Error("Expected error when storage fails, got nil")
	}

	// Verifica se o storage foi chamado
	if len(storage.GetCallHistory()) != 1 {
		t.Errorf("Expected 1 call to DeleteDomain, got %d", len(storage.GetCallHistory()))
	}
}
