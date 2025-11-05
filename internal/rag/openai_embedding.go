package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"awesomeProject2/internal/models"
)

// OpenAIEmbeddingService 实现基于OpenAI的嵌入服务
type OpenAIEmbeddingService struct {
	apiKey string
	model  string
	http   *http.Client
}

// NewOpenAIEmbeddingService 创建新的OpenAI嵌入服务
func NewOpenAIEmbeddingService(apiKey, model string) EmbeddingService {
	return &OpenAIEmbeddingService{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

// EmbeddingRequest 表示嵌入请求
type EmbeddingRequest1 struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingResponse 表示嵌入响应
type EmbeddingResponse1 struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// GenerateEmbedding 生成文本的向量表示
func (o *OpenAIEmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	// 准备请求
	reqBody := EmbeddingRequest1{
		Model: o.model,
		Input: []string{text},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	// 发送请求
	resp, err := o.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var embeddingResp EmbeddingResponse1
	err = json.Unmarshal(body, &embeddingResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned from API")
	}

	return embeddingResp.Data[0].Embedding, nil
}

// GenerateTableEmbedding 生成表结构的嵌入向量
func (o *OpenAIEmbeddingService) GenerateTableEmbedding(table *models.Table) ([]float32, error) {
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

	return o.GenerateEmbedding(tableText)
}
