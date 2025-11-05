// Package rag provides interfaces for Retrieval-Augmented Generation functionality
package rag

import (
	"awesomeProject2/internal/models"
)

// VectorStore 定义向量存储接口
type VectorStore interface {
	IndexTableStructure(table *models.Table, vector []float32) error
	SearchSimilarTables(queryVector []float32, topK int) ([]*models.Table, error)
	DeleteTableVectors(tableName string) error
}

// EmbeddingService 定义嵌入服务接口
type EmbeddingService interface {
	GenerateEmbedding(text string) ([]float32, error)
	GenerateTableEmbedding(table *models.Table) ([]float32, error)
}