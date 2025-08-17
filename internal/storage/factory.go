package storage

import (
	"github.com/luizhreis/domain-watcher/internal/storage/memory"
)

type StorageType string

const (
	StorageTypeMemory     StorageType = "memory"
	StorageTypePostgreSQL StorageType = "postgresql"
)

// StorageConfig defines the configuration for the storage system.
type StorageConfig struct {
	Type    StorageType            `json:"type"`
	Path    string                 `json:"path,omitempty"`    // Para SQLite: caminho do arquivo
	DSN     string                 `json:"dsn,omitempty"`     // Para PostgreSQL: string de conexão
	Options map[string]interface{} `json:"options,omitempty"` // Configurações específicas
}

var _ Storage = (*memory.MemoryStorage)(nil)

// NewStorage creates a new storage instance based on the provided type.
func NewStorage(config *StorageConfig) (Storage, error) {
	switch config.Type {
	case StorageTypeMemory:
		return memory.NewMemoryStorage(), nil
	case StorageTypePostgreSQL:
		return nil, ErrUnsupportedStorageType
	default:
		return nil, ErrUnsupportedStorageType
	}
}
