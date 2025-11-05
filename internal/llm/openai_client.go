package llm

import (
	"context"
	"fmt"
	"strings"

	"sql_generator/internal/config"
	"sql_generator/internal/models"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements Client for OpenAI API
type OpenAIClient struct {
	client *openai.Client
	config config.LLMConfig
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config config.LLMConfig) Client {
	client := openai.NewClient(config.APIKey)
	return &OpenAIClient{
		client: client,
		config: config,
	}
}

// GenerateSQL generates SQL using OpenAI API
func (o *OpenAIClient) GenerateSQL(description string, tables []*models.Table) (string, error) {
	// Build prompt
	fmt.Println("----openai---")
	prompt := o.buildPrompt(description, tables)

	// Prepare request
	req := openai.ChatCompletionRequest{
		Model: o.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   o.config.MaxTokens,
		Temperature: float32(o.config.Temp),
	}

	// Send request
	resp, err := o.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to generate SQL: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI API")
	}

	return resp.Choices[0].Message.Content, nil
}

// buildPrompt constructs the prompt for the LLM
func (o *OpenAIClient) buildPrompt(description string, tables []*models.Table) string {
	var prompt strings.Builder

	prompt.WriteString("根据以下表结构和用户需求生成SQL查询语句：\n\n")
	prompt.WriteString(fmt.Sprintf("用户需求：%s\n\n", description))
	prompt.WriteString("相关表结构：\n")

	// Add table structures, but control total length
	totalLength := 0
	for _, table := range tables {
		tableInfo := fmt.Sprintf("\n表名: %s\n描述: %s\n", table.Name, table.Description)

		// Check if we exceed length limit (roughly 8000 chars to leave room for rest of prompt)
		if totalLength+len(tableInfo) > 8000 {
			prompt.WriteString("\n... (为控制提示词长度，省略了部分表结构) ...\n")
			break
		}

		prompt.WriteString(tableInfo)
		prompt.WriteString("字段:\n")

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
				prompt.WriteString("  ... (字段信息省略)\n")
				break
			}

			prompt.WriteString(columnInfo)
			totalLength += len(columnInfo)
		}

		totalLength += len(tableInfo)
	}

	prompt.WriteString("\n请根据用户需求和提供的表结构生成相应的SQL查询语句：\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("1. 只返回有效的SQL语句\n")
	prompt.WriteString("2. 如需要多表关联，请使用适当的JOIN语句\n")
	prompt.WriteString("3. 不要包含任何解释或其他文本，只返回SQL\n")
	prompt.WriteString("4. 使用标准SQL语法\n")

	return prompt.String()
}
