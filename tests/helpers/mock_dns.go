package helpers

import (
	"github.com/luizhreis/domain-watcher/internal/dns"
)

// MockDNS - Mock centralizado para testes de integração
type MockDNS struct {
	resolveFunc func(domain string) (string, error)
	callHistory []string
}

// Compile-time check para garantir que implementa a interface
var _ dns.DNS = (*MockDNS)(nil)

func NewMockDNS() *MockDNS {
	return &MockDNS{
		callHistory: make([]string, 0),
	}
}

func (m *MockDNS) Resolve(domain string) (string, error) {
	m.callHistory = append(m.callHistory, domain)

	if m.resolveFunc != nil {
		return m.resolveFunc(domain)
	}
	return "192.168.1.1", nil
}

func (m *MockDNS) SetResolveFunc(f func(domain string) (string, error)) {
	m.resolveFunc = f
}

func (m *MockDNS) GetCallHistory() []string {
	return m.callHistory
}

func (m *MockDNS) ClearHistory() {
	m.callHistory = make([]string, 0)
}

func (m *MockDNS) CalledWith(domain string) bool {
	for _, call := range m.callHistory {
		if call == domain {
			return true
		}
	}
	return false
}

func (m *MockDNS) CallCount() int {
	return len(m.callHistory)
}
