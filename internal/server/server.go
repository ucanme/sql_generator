package server

import (
	"fmt"
	"net/http"
	"time"

	"awesomeProject2/internal/config"
	"awesomeProject2/internal/handlers"
	"awesomeProject2/internal/storage"
	"awesomeProject2/internal/llm"
	"awesomeProject2/internal/rag"

	"github.com/gin-gonic/gin"
)

// New creates a new HTTP server with configured routes
func New(cfg *config.Config) (*http.Server, error) {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Create MySQL storage instead of MongoDB storage
	mysqlStore, err := storage.NewMySQLStore(cfg.MySQL.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL store: %w", err)
	}

	// Create embedding service
	fmt.Printf("Creating embedding service with provider: %s\n", cfg.Embedding.Provider)
	embeddingSvc, err := rag.NewEmbeddingService(cfg.Embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding service: %w", err)
	}

	// Create vector store
	vectorStore, err := rag.NewPineconeVectorStore(cfg.VectorDB.APIKey, cfg.VectorDB.IndexName, mysqlStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	// Create RAG enhanced storage
	store := storage.NewRAGEnhancedStore(mysqlStore, embeddingSvc, vectorStore)

	// Load existing tables from MySQL and index them for RAG
	err = loadAndIndexTables(mysqlStore, embeddingSvc, vectorStore)
	if err != nil {
		fmt.Printf("Warning: failed to load and index tables: %v\n", err)
	}

	// Create LLM client with RAG enhancement
	var baseLLMClient llm.Client
	// 根据配置的模型名称来选择合适的LLM客户端
	switch {
	case cfg.LLM.Model == "deepseek" || cfg.LLM.Model == "deepseek-chat":
		baseLLMClient = llm.NewDeepSeekClient(cfg.LLM)
	case cfg.LLM.Model == "qwen" || cfg.LLM.Model == "qwen-turbo" || cfg.LLM.Model == "qwen-plus" || cfg.LLM.Model == "qwen-max":
		baseLLMClient = llm.NewDeepSeekClient(cfg.LLM) // 使用DeepSeek客户端，因为它与Qwen兼容
	default:
		baseLLMClient = llm.NewOpenAIClient(cfg.LLM)
	}

	llmClient := llm.NewRAGEnhancedClient(cfg.LLM, baseLLMClient, embeddingSvc, vectorStore)

	// Create handlers
	handler := handlers.NewHandler(store, llmClient)

	// Register routes
	handler.RegisterRoutes(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return srv, nil
}

// loadAndIndexTables loads existing tables from storage and indexes them for RAG
func loadAndIndexTables(store storage.Store, embeddingSvc rag.EmbeddingService, vectorStore rag.VectorStore) error {
	fmt.Println("Loading existing tables from MySQL...")

	// Get all tables
	tables, err := store.ListTables(1000, 0) // Load up to 1000 tables
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	fmt.Printf("Loaded %d tables from MySQL\n", len(tables))

	// Index each table for RAG
	for _, table := range tables {
		fmt.Printf("Indexing table: %s (%s)\n", table.Name, table.Description)
		
		// Generate embedding for the table
		vector, err := embeddingSvc.GenerateTableEmbedding(table)
		if err != nil {
			fmt.Printf("Warning: failed to generate embedding for table %s: %v\n", table.Name, err)
			continue
		}

		// Index the table structure in vector store
		err = vectorStore.IndexTableStructure(table, vector)
		if err != nil {
			fmt.Printf("Warning: failed to index table %s: %v\n", table.Name, err)
			continue
		}
		
		fmt.Printf("Successfully indexed table: %s\n", table.Name)
	}

	return nil
}