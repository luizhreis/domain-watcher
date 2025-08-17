package helpers

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luizhreis/domain-watcher/internal/models"
)

// TestDomainBuilder - Builder pattern para criar domains de teste
type TestDomainBuilder struct {
	domain *models.Domain
}

func NewTestDomainBuilder() *TestDomainBuilder {
	return &TestDomainBuilder{
		domain: &models.Domain{
			ID:   uuid.New(),
			Name: "Test Domain",
			URL:  "test.example.com",
		},
	}
}

func (b *TestDomainBuilder) WithID(id uuid.UUID) *TestDomainBuilder {
	b.domain.ID = id
	return b
}

func (b *TestDomainBuilder) WithName(name string) *TestDomainBuilder {
	b.domain.Name = name
	return b
}

func (b *TestDomainBuilder) WithURL(url string) *TestDomainBuilder {
	b.domain.URL = url
	return b
}

func (b *TestDomainBuilder) Build() *models.Domain {
	return b.domain
}

// TestCheckResultBuilder - Builder pattern para criar check results
type TestCheckResultBuilder struct {
	result *models.CheckResult
}

func NewTestCheckResultBuilder() *TestCheckResultBuilder {
	return &TestCheckResultBuilder{
		result: &models.CheckResult{
			ID:            uuid.New(),
			DomainID:      uuid.New(),
			StatusCode:    200,
			ResponseTime:  100,
			Error:         "",
			RedirectURL:   "",
			RedirectCount: 0,
			CheckedAt:     time.Now(),
			ContentLength: 1024,
			Server:        "nginx",
			ResolvedIP:    "192.168.1.1",
		},
	}
}

func (b *TestCheckResultBuilder) WithID(id uuid.UUID) *TestCheckResultBuilder {
	b.result.ID = id
	return b
}

func (b *TestCheckResultBuilder) WithDomainID(domainID uuid.UUID) *TestCheckResultBuilder {
	b.result.DomainID = domainID
	return b
}

func (b *TestCheckResultBuilder) WithStatusCode(statusCode int) *TestCheckResultBuilder {
	b.result.StatusCode = statusCode
	return b
}

func (b *TestCheckResultBuilder) WithError(errorMsg string) *TestCheckResultBuilder {
	b.result.Error = errorMsg
	return b
}

func (b *TestCheckResultBuilder) WithResolvedIP(ip string) *TestCheckResultBuilder {
	b.result.ResolvedIP = ip
	return b
}

func (b *TestCheckResultBuilder) WithResponseTime(responseTime int64) *TestCheckResultBuilder {
	b.result.ResponseTime = responseTime
	return b
}

func (b *TestCheckResultBuilder) WithCheckedAt(timestamp time.Time) *TestCheckResultBuilder {
	b.result.CheckedAt = timestamp
	return b
}

func (b *TestCheckResultBuilder) Build() *models.CheckResult {
	return b.result
}

// TestCheckResultMatcher - Matcher pattern para validação
type TestCheckResultMatcher struct {
	result *models.CheckResult
	errors []string
}

func NewTestCheckResultMatcher(result *models.CheckResult) *TestCheckResultMatcher {
	return &TestCheckResultMatcher{
		result: result,
		errors: make([]string, 0),
	}
}

func (m *TestCheckResultMatcher) HasID(expectedID uuid.UUID) *TestCheckResultMatcher {
	if m.result.ID != expectedID {
		m.errors = append(m.errors,
			fmt.Sprintf("Expected ID %s, got %s", expectedID, m.result.ID))
	}
	return m
}

func (m *TestCheckResultMatcher) HasDomainID(expectedDomainID uuid.UUID) *TestCheckResultMatcher {
	if m.result.DomainID != expectedDomainID {
		m.errors = append(m.errors,
			fmt.Sprintf("Expected DomainID %s, got %s", expectedDomainID, m.result.DomainID))
	}
	return m
}

func (m *TestCheckResultMatcher) HasStatusCode(expectedStatusCode int) *TestCheckResultMatcher {
	if m.result.StatusCode != expectedStatusCode {
		m.errors = append(m.errors,
			fmt.Sprintf("Expected StatusCode %d, got %d", expectedStatusCode, m.result.StatusCode))
	}
	return m
}

func (m *TestCheckResultMatcher) HasResolvedIP(expectedIP string) *TestCheckResultMatcher {
	if m.result.ResolvedIP != expectedIP {
		m.errors = append(m.errors,
			fmt.Sprintf("Expected ResolvedIP %s, got %s", expectedIP, m.result.ResolvedIP))
	}
	return m
}

func (m *TestCheckResultMatcher) HasNoError() *TestCheckResultMatcher {
	if m.result.Error != "" {
		m.errors = append(m.errors,
			fmt.Sprintf("Expected no error, got %s", m.result.Error))
	}
	return m
}

func (m *TestCheckResultMatcher) HasValidTimestamp() *TestCheckResultMatcher {
	if m.result.CheckedAt.IsZero() {
		m.errors = append(m.errors, "Expected non-zero timestamp")
	}
	return m
}

func (m *TestCheckResultMatcher) GetErrors() []string {
	return m.errors
}

func (m *TestCheckResultMatcher) IsValid() bool {
	return len(m.errors) == 0
}

// Helper para assertiva
func AssertValidCheckResult(t testing.TB, result *models.CheckResult) {
	if result == nil {
		t.Fatal("CheckResult should not be nil")
	}

	matcher := NewTestCheckResultMatcher(result)
	matcher.HasValidTimestamp()

	if !matcher.IsValid() {
		for _, err := range matcher.GetErrors() {
			t.Error(err)
		}
	}
}
