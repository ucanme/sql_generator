package llm

import (
	"fmt"

	"sql_generator/internal/config"
	"sql_generator/internal/models"
	"sql_generator/internal/rag"
)

// RAGEnhancedClient 实现结合RAG的LLM客户端
type RAGEnhancedClient struct {
	config       config.LLMConfig
	baseClient   Client
	embeddingSvc rag.EmbeddingService
	vectorStore  rag.VectorStore
}

// NewRAGEnhancedClient 创建新的RAG增强客户端
func NewRAGEnhancedClient(config config.LLMConfig, baseClient Client, embeddingSvc rag.EmbeddingService, vectorStore rag.VectorStore) Client {
	return &RAGEnhancedClient{
		config:       config,
		baseClient:   baseClient,
		embeddingSvc: embeddingSvc,
		vectorStore:  vectorStore,
	}
}

// GenerateSQL 使用RAG增强的方式生成SQL
func (r *RAGEnhancedClient) GenerateSQL(description string, tables []*models.Table) (string, error) {
	// 如果没有提供表，则使用RAG检索相关表
	if len(tables) == 0 {
		var err error
		tables, err = r.retrieveRelevantTables(description)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve relevant tables: %w", err)
		}
	}

	fmt.Println(r.config.Model)
	// 使用基础客户端生成SQL
	return r.baseClient.GenerateSQL(description, tables)
}

// retrieveRelevantTables 使用RAG检索相关表
func (r *RAGEnhancedClient) retrieveRelevantTables(description string) ([]*models.Table, error) {
	// 生成查询的向量表示
	queryVector, err := r.embeddingSvc.GenerateEmbedding(description)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// 在向量数据库中搜索相似表
	tables, err := r.vectorStore.SearchSimilarTables(queryVector, 10) // 获取前10个最相关表
	if err != nil {
		return nil, fmt.Errorf("failed to search similar tables: %w", err)
	}

	return tables, nil
}

// buildPrompt 构建提示词
func (r *RAGEnhancedClient) buildPrompt(description string, tables []*models.Table) string {
	prompt := "根据以下表结构和用户需求生成SQL查询语句：\n\n"
	prompt += fmt.Sprintf("用户需求：%s\n\n", description)
	prompt += "相关表结构：\n"

	// 添加表结构信息，但控制总长度
	totalLength := 0
	for _, table := range tables {
		tableInfo := fmt.Sprintf("\n表名: %s\n描述: %s\n", table.Name, table.Description)

		// 检查是否超出长度限制（大约8000字符）
		if totalLength+len(tableInfo) > 8000 {
			prompt += "\n... (为控制提示词长度，省略了部分表结构) ...\n"
			break
		}

		prompt += tableInfo
		prompt += "字段:\n"

		for _, column := range table.Columns {
			columnInfo := fmt.Sprintf("  - %s (%s): %s", column.Name, column.Type, column.Description)
			if column.IsPrimary {
				columnInfo += " [主键]"
			}
			if column.IsRequired {
				columnInfo += " [必填]"
			}
			columnInfo += "\n"

			if totalLength+len(columnInfo) > 8000 {
				prompt += "  ... (字段信息省略)\n"
				break
			}

			prompt += columnInfo
			totalLength += len(columnInfo)
		}

		totalLength += len(tableInfo)
	}

	prompt += "\n请根据用户需求和提供的表结构生成相应的SQL查询语句：\n"
	prompt += "要求：\n"
	prompt += "1. 只返回有效的SQL语句\n"
	prompt += "2. 如需要多表关联，请使用适当的JOIN语句\n"
	prompt += "3. 不要包含任何解释或其他文本，只返回SQL\n"
	prompt += "4. 使用标准SQL语法\n"

	return prompt
}
