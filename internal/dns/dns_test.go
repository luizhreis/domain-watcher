package dns

import (
	"net"
	"testing"
)

// TestNewDNS testa a criação do resolver (white-box)
func TestNewDNS(t *testing.T) {
	resolver := NewDNS()

	if resolver == nil {
		t.Error("NewDNS returned nil")
	}
}

// TestDNSResolve testa a resolução DNS (white-box)
func TestDNSResolve(t *testing.T) {
	resolver := NewDNS()

	t.Run("Valid Domain", func(t *testing.T) {
		// Teste com domínio real
		ip, err := resolver.Resolve("google.com")

		if err != nil {
			t.Errorf("Expected no error for google.com, got %v", err)
		}

		if ip == "" {
			t.Error("Expected non-empty IP address")
		}

		// Valida formato IP
		if net.ParseIP(ip) == nil {
			t.Errorf("Expected valid IP address, got %s", ip)
		}
	})

	t.Run("Invalid Domain", func(t *testing.T) {
		// Teste com domínio vazio
		ip, err := resolver.Resolve("")

		if err == nil {
			t.Error("Expected error for empty domain")
		}

		if ip != "" {
			t.Errorf("Expected empty IP for invalid domain, got %s", ip)
		}
	})
}
