package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"awesomeProject2/internal/config"
	"awesomeProject2/internal/models"
)

// DeepSeekClient implements Client for DeepSeek API
type DeepSeekClient struct {
	config config.LLMConfig
	http   *http.Client
}

// NewDeepSeekClient creates a new DeepSeek client
func NewDeepSeekClient(config config.LLMConfig) Client {
	return &DeepSeekClient{
		config: config,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatCompletionRequest represents the request to DeepSeek Chat API
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response from DeepSeek Chat API
type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Choices []ChatChoice `json:"choices"`
}

// ChatChoice represents a choice in the chat completion
type ChatChoice struct {
	Message ChatMessage `json:"message"`
}

// GenerateSQL generates SQL using DeepSeek API
func (c *DeepSeekClient) GenerateSQL(description string, tables []*models.Table) (string, error) {
	// Build prompt
	prompt := c.buildPrompt(description, tables)

	// Prepare request
	reqBody := ChatCompletionRequest{
		Model: c.config.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temp,
	}

	jsonData, err := json.Marshal(reqBody)

	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request - 支持多种模型提供商
	var apiURL string
	if c.config.Model == "qwen" || c.config.Model == "qwen-turbo" || c.config.Model == "qwen-plus" || c.config.Model == "qwen-max" {
		// 阿里云千问API
		apiURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	} else {
		// DeepSeek API
		apiURL = "https://api.deepseek.com/v1/chat/completions"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// Send request
	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var completionResp ChatCompletionResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return completionResp.Choices[0].Message.Content, nil
}

// buildPrompt constructs the prompt for the LLM
func (c *DeepSeekClient) buildPrompt(description string, tables []*models.Table) string {
	println("-----proto---start--")
	prompt := "根据以下表结构和用户需求生成SQL查询语句：\n\n"
	prompt += fmt.Sprintf("用户需求：%s\n\n", description)
	prompt += "相关表结构：\n"

	// Add table structures, but control total length
	totalLength := 0
	fmt.Println("---len(tables)--", len(tables))
	for _, table := range tables {
		tableInfo := fmt.Sprintf("\n表名: %s\n描述: %s\n", table.Name, table.Description)

		// Check if we exceed length limit (roughly 8000 chars to leave room for rest of prompt)
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
