package storage

import (
	"fmt"

	"sql_generator/internal/models"
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

// RAGEnhancedStore 结合RAG功能的存储实现
type RAGEnhancedStore struct {
	Store        // 修改这里，使用通用的Store接口而不是*MongoStore
	embeddingSvc EmbeddingService
	vectorStore  VectorStore
}

// NewRAGEnhancedStore 创建支持RAG的存储实例
func NewRAGEnhancedStore(store Store, embeddingSvc EmbeddingService, vectorStore VectorStore) Store {
	return &RAGEnhancedStore{
		Store:        store,
		embeddingSvc: embeddingSvc,
		vectorStore:  vectorStore,
	}
}

// CreateTable 保存表定义并创建向量索引
func (r *RAGEnhancedStore) CreateTable(table *models.Table) error {
	// 首先保存到存储中
	err := r.Store.CreateTable(table)
	if err != nil {
		return fmt.Errorf("failed to create table in storage: %w", err)
	}

	// 生成并向量数据库中索引表结构
	err = r.indexTableForRAG(table)
	if err != nil {
		// 如果向量索引失败，记录日志但不中断操作
		fmt.Printf("Warning: Failed to index table for RAG: %v\n", err)
	}

	return nil
}

// UpdateTable 更新表并更新向量索引
func (r *RAGEnhancedStore) UpdateTable(name string, table *models.Table) error {
	// 更新存储中的表
	err := r.Store.UpdateTable(name, table)
	if err != nil {
		return fmt.Errorf("failed to update table in storage: %w", err)
	}

	// 删除旧的向量索引
	err = r.vectorStore.DeleteTableVectors(name)
	if err != nil {
		fmt.Printf("Warning: Failed to delete old vector index: %v\n", err)
	}

	// 创建新的向量索引
	err = r.indexTableForRAG(table)
	if err != nil {
		fmt.Printf("Warning: Failed to reindex table for RAG: %v\n", err)
	}

	return nil
}

// DeleteTable 删除表并删除向量索引
func (r *RAGEnhancedStore) DeleteTable(name string) error {
	// 从存储中删除表
	err := r.Store.DeleteTable(name)
	if err != nil {
		return fmt.Errorf("failed to delete table from storage: %w", err)
	}

	// 删除向量索引
	err = r.vectorStore.DeleteTableVectors(name)
	if err != nil {
		fmt.Printf("Warning: Failed to delete vector index: %v\n", err)
	}

	return nil
}

// indexTableForRAG 为表创建向量索引
func (r *RAGEnhancedStore) indexTableForRAG(table *models.Table) error {
	// 生成表结构的向量表示
	vector, err := r.embeddingSvc.GenerateTableEmbedding(table)
	if err != nil {
		return fmt.Errorf("failed to generate table embedding: %w", err)
	}

	// 索引到向量数据库
	err = r.vectorStore.IndexTableStructure(table, vector)
	if err != nil {
		return fmt.Errorf("failed to index table structure: %w", err)
	}

	return nil
}

// SearchTables 搜索表（结合传统搜索和RAG向量搜索）
func (r *RAGEnhancedStore) SearchTables(keyword string, limit, offset int) ([]*models.Table, error) {
	// 优先尝试使用RAG向量搜索
	tables, err := r.searchTablesByEmbedding(keyword, limit)
	if err == nil && len(tables) > 0 {
		return tables, nil
	}

	// 如果RAG搜索失败或未找到结果，回退到传统的文本搜索
	return r.Store.SearchTables(keyword, limit, offset)
}

// searchTablesByEmbedding 使用嵌入向量搜索表
func (r *RAGEnhancedStore) searchTablesByEmbedding(keyword string, limit int) ([]*models.Table, error) {
	// 生成查询关键词的向量表示
	queryVector, err := r.embeddingSvc.GenerateEmbedding(keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// 在向量数据库中搜索
	tables, err := r.vectorStore.SearchSimilarTables(queryVector, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search by embedding: %w", err)
	}

	return tables, nil
}
