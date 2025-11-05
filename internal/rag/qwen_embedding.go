package rag

import (
	"awesomeProject2/internal/models"
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// QwenEmbeddingService 使用OpenAI SDK实现的阿里云千问嵌入服务
type QwenEmbeddingService struct {
	client *openai.Client
	model  string
}

// NewQwenEmbeddingService 创建新的千问嵌入服务（使用OpenAI SDK）
func NewQwenEmbeddingService(apiKey, model string) EmbeddingService {
	// 配置阿里云千问的API端点
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"

	// 创建OpenAI客户端
	client := openai.NewClientWithConfig(config)

	return &QwenEmbeddingService{
		client: client,
		model:  model,
	}
}

// GenerateEmbedding 生成文本的向量表示
func (q *QwenEmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	// 创建嵌入请求
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(q.model),
	}

	// 发送请求
	resp, err := q.client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned from API")
	}

	return resp.Data[0].Embedding, nil
}

// GenerateTableEmbedding 生成表结构的嵌入向量
func (q *QwenEmbeddingService) GenerateTableEmbedding(table *models.Table) ([]float32, error) {
	// 构造表结构的文本表示
	tableText := fmt.Sprintf("Table name: %s\nDescription: %s\n", table.Name, table.Description)

	for _, column := range table.Columns {
		columnText := fmt.Sprintf("Column: %s, Type: %s, Description: %s",
			column.Name, column.Type, column.Description)
		if column.IsPrimary {
			columnText += ", Primary Key"
		}
		if column.IsRequired {
			columnText += ", Required"
		}
		tableText += columnText + "\n"
	}

	return q.GenerateEmbedding(tableText)
}