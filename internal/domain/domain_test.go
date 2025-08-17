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