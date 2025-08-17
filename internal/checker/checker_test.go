package checker

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
)

// MockDNS é um mock interno para testes unitários
type MockDNS struct {
	resolveFunc func(domain string) (string, error)
}

func (m *MockDNS) Resolve(domain string) (string, error) {
	if m.resolveFunc != nil {
		return m.resolveFunc(domain)
	}
	return "192.168.1.1", nil
}

// TestNewChecker testa a criação do checker (white-box)
func TestNewChecker(t *testing.T) {
	mockDNS := &MockDNS{}

	checkerInstance := NewChecker(mockDNS)

	if checkerInstance == nil {
		t.Error("NewChecker returned nil")
	}

	// Acesso interno - verifica se o campo foi definido corretamente
	if c, ok := checkerInstance.(*checker); ok {
		if c.DNSResolver != mockDNS {
			t.Error("DNSResolver field not set correctly")
		}
	} else {
		t.Error("NewChecker did not return correct internal type")
	}
}

// TestCheckerCheckDomain testa CheckDomain com acesso interno (white-box)
func TestCheckerCheckDomain(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDNS := &MockDNS{}
		expectedIP := "192.168.1.200"

		mockDNS.resolveFunc = func(domain string) (string, error) {
			return expectedIP, nil
		}

		checkerInstance := NewChecker(mockDNS)
		testDomainID := uuid.New()
		domain := &models.Domain{
			ID:   testDomainID,
			Name: "Unit Test",
			URL:  "unit.test",
		}

		beforeTime := time.Now()
		result, err := checkerInstance.CheckDomain(domain)
		afterTime := time.Now()

		// Testes unitários - verifica TODA a implementação
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("Expected result to not be nil")
		}

		// White-box: testa valores específicos da implementação
		if result.ID == uuid.Nil {
			t.Error("Expected ID to be generated, got nil UUID")
		}

		if result.DomainID != domain.ID {
			t.Errorf("Expected DomainID %d, got %d", domain.ID, result.DomainID)
		}

		if result.StatusCode != 200 {
			t.Errorf("Expected StatusCode 200, got %d", result.StatusCode)
		}

		if result.ResponseTime != 100 {
			t.Errorf("Expected ResponseTime 100, got %d", result.ResponseTime)
		}

		if result.Error != "" {
			t.Errorf("Expected empty error, got %s", result.Error)
		}

		if result.RedirectURL != "" {
			t.Errorf("Expected empty RedirectURL, got %s", result.RedirectURL)
		}

		if result.RedirectCount != 0 {
			t.Errorf("Expected RedirectCount 0, got %d", result.RedirectCount)
		}

		if result.ContentLength != 1024 {
			t.Errorf("Expected ContentLength 1024, got %d", result.ContentLength)
		}

		if result.Server != "nginx" {
			t.Errorf("Expected Server 'nginx', got %s", result.Server)
		}

		if result.ResolvedIP != expectedIP {
			t.Errorf("Expected ResolvedIP %s, got %s", expectedIP, result.ResolvedIP)
		}

		// White-box: testa timestamp interno detalhadamente
		if result.CheckedAt.IsZero() {
			t.Error("CheckedAt should not be zero")
		}

		if result.CheckedAt.Before(beforeTime) || result.CheckedAt.After(afterTime) {
			t.Errorf("CheckedAt timestamp %v should be between %v and %v",
				result.CheckedAt, beforeTime, afterTime)
		}
	})

	t.Run("DNS Error", func(t *testing.T) {
		mockDNS := &MockDNS{}
		expectedError := errors.New("DNS lookup failed")

		mockDNS.resolveFunc = func(domain string) (string, error) {
			return "", expectedError
		}

		checkerInstance := NewChecker(mockDNS)
		testDomainID := uuid.New()
		domain := &models.Domain{
			ID:   testDomainID,
			Name: "Failed Domain",
			URL:  "failed.test",
		}

		result, err := checkerInstance.CheckDomain(domain)

		// White-box: testa propagação de erro específica
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err != expectedError {
			t.Errorf("Expected exact error reference, got different error")
		}

		if result != nil {
			t.Error("Expected result to be nil when error occurs")
		}
	})

	t.Run("Nil Domain Panic", func(t *testing.T) {
		mockDNS := &MockDNS{}
		checkerInstance := NewChecker(mockDNS)

		// White-box: testa comportamento interno específico
		defer func() {
			if r := recover(); r != nil {
				t.Logf("CheckDomain panicked as expected: %v", r)
			} else {
				t.Error("Expected panic with nil domain")
			}
		}()

		_, _ = checkerInstance.CheckDomain(nil)
	})
}

// BenchmarkCheckerCheckDomain benchmark unitário
func BenchmarkCheckerCheckDomain(b *testing.B) {
	mockDNS := &MockDNS{}
	mockDNS.resolveFunc = func(domain string) (string, error) {
		return "192.168.1.1", nil
	}

	checkerInstance := NewChecker(mockDNS)
	testDomainID := uuid.New()
	domain := &models.Domain{
		ID:   testDomainID,
		Name: "Benchmark Domain",
		URL:  "benchmark.test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := checkerInstance.CheckDomain(domain)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
