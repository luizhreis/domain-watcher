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
