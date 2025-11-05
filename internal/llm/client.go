package llm

import (
	"sql_generator/internal/models"
)

// Client defines the interface for LLM clients
type Client interface {
	GenerateSQL(description string, tables []*models.Table) (string, error)
}
