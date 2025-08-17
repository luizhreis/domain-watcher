package integration

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/checker"
	"github.com/luizhreis/domain-watcher/tests/helpers"
)

// TestCheckerIntegration - Testes black-box de integração
func TestCheckerIntegration(t *testing.T) {
	t.Run("End-to-End Success Flow", func(t *testing.T) {
		// Arrange - integração real entre componentes
		mockDNS := helpers.NewMockDNS()
		expectedIP := "203.0.113.1"
		testDomainID := uuid.New()

		mockDNS.SetResolveFunc(func(domain string) (string, error) {
			switch domain {
			case "production.example.com":
				return expectedIP, nil
			case "staging.example.com":
				return "203.0.113.2", nil
			default:
				return "", errors.New("domain not found")
			}
		})

		checkerInstance := checker.NewChecker(mockDNS)
		domain := helpers.NewTestDomainBuilder().
			WithID(testDomainID).
			WithName("Production Domain").
			WithURL("production.example.com").
			Build()

		// Act - comportamento de ponta a ponta
		result, err := checkerInstance.CheckDomain(domain)

		// Assert - validação black-box do comportamento
		if err != nil {
			t.Fatalf("Integration test failed with error: %v", err)
		}

		// Usa helpers para validação fluente
		matcher := helpers.NewTestCheckResultMatcher(result).
			HasDomainID(testDomainID).
			HasResolvedIP(expectedIP).
			HasStatusCode(200).
			HasNoError().
			HasValidTimestamp()

		if !matcher.IsValid() {
			for _, err := range matcher.GetErrors() {
				t.Error(err)
			}
		}

		// Verifica interação com DNS
		if !mockDNS.CalledWith("production.example.com") {
			t.Error("DNS should have been called with production.example.com")
		}

		if mockDNS.CallCount() != 1 {
			t.Errorf("Expected 1 DNS call, got %d", mockDNS.CallCount())
		}
	})

	t.Run("End-to-End Error Flow", func(t *testing.T) {
		// Arrange - cenário de erro integrado
		mockDNS := helpers.NewMockDNS()
		dnsError := errors.New("network timeout")

		mockDNS.SetResolveFunc(func(domain string) (string, error) {
			return "", dnsError
		})

		checkerInstance := checker.NewChecker(mockDNS)
		domain := helpers.NewTestDomainBuilder().
			WithURL("unreachable.example.com").
			Build()

		// Act
		result, err := checkerInstance.CheckDomain(domain)

		// Assert - comportamento de erro integrado
		if err == nil {
			t.Error("Expected error but got none")
		}

		if err != dnsError {
			t.Errorf("Expected specific DNS error, got %v", err)
		}

		if result != nil {
			t.Error("Result should be nil on error")
		}

		// Verifica que DNS foi chamado mesmo com erro
		if !mockDNS.CalledWith("unreachable.example.com") {
			t.Error("DNS should have been called even on error")
		}
	})

	t.Run("Multiple Domain Integration", func(t *testing.T) {
		// Arrange - teste de múltiplos domínios
		mockDNS := helpers.NewMockDNS()

		mockDNS.SetResolveFunc(func(domain string) (string, error) {
			switch domain {
			case "api.example.com":
				return "10.0.0.1", nil
			case "web.example.com":
				return "10.0.0.2", nil
			case "cdn.example.com":
				return "10.0.0.3", nil
			default:
				return "", errors.New("unknown domain")
			}
		})

		checkerInstance := checker.NewChecker(mockDNS)

		// Cria UUIDs únicos para cada domínio
		domainID1 := uuid.New()
		domainID2 := uuid.New()
		domainID3 := uuid.New()

		domains := []*helpers.TestDomainBuilder{
			helpers.NewTestDomainBuilder().WithID(domainID1).WithURL("api.example.com"),
			helpers.NewTestDomainBuilder().WithID(domainID2).WithURL("web.example.com"),
			helpers.NewTestDomainBuilder().WithID(domainID3).WithURL("cdn.example.com"),
		}

		expectedIPs := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}

		// Act - verifica integração com múltiplas chamadas
		for i, domainBuilder := range domains {
			domain := domainBuilder.Build()
			result, err := checkerInstance.CheckDomain(domain)

			// Assert - cada resultado individualmente
			if err != nil {
				t.Errorf("Domain %d failed: %v", i+1, err)
				continue
			}

			matcher := helpers.NewTestCheckResultMatcher(result).
				HasDomainID(domain.ID).
				HasResolvedIP(expectedIPs[i]).
				HasStatusCode(200)

			if !matcher.IsValid() {
				for _, err := range matcher.GetErrors() {
					t.Errorf("Domain %d: %s", i+1, err)
				}
			}
		}

		// Assert - verifica histórico completo
		if mockDNS.CallCount() != 3 {
			t.Errorf("Expected 3 DNS calls, got %d", mockDNS.CallCount())
		}

		expectedCalls := []string{"api.example.com", "web.example.com", "cdn.example.com"}
		for _, expectedCall := range expectedCalls {
			if !mockDNS.CalledWith(expectedCall) {
				t.Errorf("DNS should have been called with %s", expectedCall)
			}
		}
	})
}

// BenchmarkCheckerIntegration - benchmark de integração
func BenchmarkCheckerIntegration(b *testing.B) {
	mockDNS := helpers.NewMockDNS()
	mockDNS.SetResolveFunc(func(domain string) (string, error) {
		return "192.168.1.1", nil
	})

	checkerInstance := checker.NewChecker(mockDNS)
	domain := helpers.NewTestDomainBuilder().
		WithURL("benchmark.example.com").
		Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := checkerInstance.CheckDomain(domain)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
