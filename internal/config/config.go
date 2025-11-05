package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server    ServerConfig
	Mongo     MongoConfig
	MySQL     MySQLConfig  // 添加MySQL配置
	LLM       LLMConfig
	Embedding EmbeddingConfig
	VectorDB  VectorDBConfig
}

// ServerConfig holds the HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// MongoConfig holds the MongoDB configuration
type MongoConfig struct {
	URI      string
	Database string
}

// MySQLConfig holds the MySQL configuration
type MySQLConfig struct {
	DSN      string
	Database string
}

// LLMConfig holds the Large Language Model configuration
type LLMConfig struct {
	APIKey    string
	Model     string
	MaxTokens int
	Temp      float64
}

// EmbeddingConfig holds the embedding service configuration
type EmbeddingConfig struct {
	APIKey string
	Model  string
	// 添加Provider字段来指定使用哪个提供商
	Provider string
	// Hugging Face特定配置
	HFEndpoint string
	HFModel   string
	// Qwen特定配置
	QwenModel string
}

// VectorDBConfig holds the vector database configuration
type VectorDBConfig struct {
	APIKey      string
	IndexName   string
	Environment string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30),
		},
		Mongo: MongoConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "sqlbot"),
		},
		MySQL: MySQLConfig{
			DSN:      getEnv("MYSQL_DSN", "root:password@tcp(localhost:3306)/sqlbot?charset=utf8mb4&parseTime=True&loc=Local"),
			Database: getEnv("MYSQL_DATABASE", "sqlbot"),
		},
		LLM: LLMConfig{
			APIKey:    getEnv("LLM_API_KEY", ""),
			Model:     getEnv("LLM_MODEL", "gpt-3.5-turbo"),
			MaxTokens: getEnvAsInt("LLM_MAX_TOKENS", 2000),
			Temp:      getEnvAsFloat("LLM_TEMPERATURE", 0.3),
		},
		Embedding: EmbeddingConfig{
			APIKey:    getEnv("EMBEDDING_API_KEY", ""),
			Model:     getEnv("EMBEDDING_MODEL", "text-embedding-v1"),
			Provider:  getEnv("EMBEDDING_PROVIDER", "qwen"), // 默认使用阿里云千问
			HFEndpoint: getEnv("HF_ENDPOINT", "https://api-inference.huggingface.co/models/"),
			HFModel:   getEnv("HF_MODEL", "sentence-transformers/all-MiniLM-L6-v2"),
			QwenModel: getEnv("QWEN_MODEL", "text-embedding-v1"),
		},
		VectorDB: VectorDBConfig{
			APIKey:      getEnv("VECTOR_DB_API_KEY", ""),
			IndexName:   getEnv("VECTOR_DB_INDEX_NAME", "sqlbot-tables"),
			Environment: getEnv("VECTOR_DB_ENVIRONMENT", "us-west1-gcp"),
		},
	}

	return cfg, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}