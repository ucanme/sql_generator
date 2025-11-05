package llm

import (
	"awesomeProject2/internal/models"
)

// Client defines the interface for LLM clients
type Client interface {
	GenerateSQL(description string, tables []*models.Table) (string, error)
}