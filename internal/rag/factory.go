package rag

import (
	"fmt"
	"sql_generator/internal/config"
)

// NewEmbeddingService 根据配置创建嵌入服务
func NewEmbeddingService(config config.EmbeddingConfig) (EmbeddingService, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIEmbeddingService(config.APIKey, config.Model), nil
	case "qwen":
		// 使用阿里云千问作为嵌入服务提供者
		// 如果配置了专门的千问模型，则使用它，否则使用通用模型配置
		model := config.QwenModel
		if model == "" {
			model = config.Model
		}
		return NewQwenEmbeddingService(config.APIKey, model), nil
	case "huggingface":
		// 使用Hugging Face作为嵌入服务提供者
		// 如果有自定义端点配置，则使用它
		if config.HFEndpoint != "" {
			return NewHFEmbeddingServiceWithConfig(config.APIKey, config.HFModel, config.HFEndpoint), nil
		}
		return NewHFEmbeddingService(config.APIKey, config.HFModel), nil
	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", config.Provider)
	}
}
