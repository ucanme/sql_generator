package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"sql_generator/internal/models"
)

// HFEmbeddingService implements Hugging Face embedding service
type HFEmbeddingService struct {
	apiKey  string
	model   string
	baseURL string
	http    *http.Client
}

// NewHFEmbeddingService creates a new Hugging Face embedding service
func NewHFEmbeddingService(apiKey, model string) EmbeddingService {
	// Create a custom transport with better connection handling
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &HFEmbeddingService{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api-inference.huggingface.co/models/",
		http: &http.Client{
			Transport: transport,
			Timeout:   60 * time.Second,
		},
	}
}

// NewHFEmbeddingServiceWithConfig creates a new Hugging Face embedding service with custom config
func NewHFEmbeddingServiceWithConfig(apiKey, model, endpoint string) EmbeddingService {
	// Use custom endpoint if provided, otherwise use default
	baseURL := "https://api-inference.huggingface.co/models/"
	if endpoint != "" {
		baseURL = strings.TrimRight(endpoint, "/") + "/"
	}

	// Create a custom transport with better connection handling
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &HFEmbeddingService{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		http: &http.Client{
			Transport: transport,
			Timeout:   60 * time.Second,
		},
	}
}

// HFEmbeddingRequest represents Hugging Face embedding request
type HFEmbeddingRequest struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

// HFEmbeddingSingleResponse represents single text Hugging Face embedding response
type HFEmbeddingSingleResponse struct {
	Embedding []float32 `json:"embedding"`
}

// HFEmbeddingMultiResponse represents multiple texts Hugging Face embedding response
type HFEmbeddingMultiResponse []HFEmbeddingSingleResponse

// HFEmbeddingErrorResponse represents Hugging Face error response
type HFEmbeddingErrorResponse struct {
	Error string `json:"error"`
}

// GenerateEmbedding generates vector representation of text
func (h *HFEmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	// Prepare request
	reqBody := HFEmbeddingRequest{
		Inputs: text,
		Parameters: map[string]interface{}{
			"pooling": "mean",
		},
		Options: map[string]interface{}{
			"wait_for_model": true,
			"use_cache":      true,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request with properly constructed URL
	url := h.baseURL + strings.TrimLeft(h.model, "/")

	// Retry logic for handling timeouts and temporary issues
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			cancel()
			lastErr = fmt.Errorf("failed to create request to %s: %w", url, err)
			continue
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		if h.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+h.apiKey)
		}

		// Send request
		resp, err := h.http.Do(req)
		if err != nil {
			cancel()
			// Check if it's a network error that might be temporary
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
				lastErr = fmt.Errorf("temporary network error: %w", err)
				continue
			}
			lastErr = fmt.Errorf("failed to send request to %s: %w", url, err)
			continue
		}

		// Read response
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		cancel()

		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		// Handle non-success status codes
		if resp.StatusCode != http.StatusOK {
			// Try to parse error response
			var errorResp HFEmbeddingErrorResponse
			if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
				// If it's a model loading error, we should retry
				if strings.Contains(errorResp.Error, "currently loading") ||
					strings.Contains(errorResp.Error, "timeout") ||
					resp.StatusCode == http.StatusInternalServerError ||
					resp.StatusCode == http.StatusTooManyRequests {
					lastErr = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, errorResp.Error)
					continue
				}
				return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, errorResp.Error)
			}
			lastErr = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			continue
		}

		// Parse response
		var singleResp HFEmbeddingSingleResponse
		var multiResp HFEmbeddingMultiResponse

		// Try to parse as single response
		if err := json.Unmarshal(body, &singleResp); err == nil && len(singleResp.Embedding) > 0 {
			return singleResp.Embedding, nil
		}

		// Try to parse as multiple response
		if err := json.Unmarshal(body, &multiResp); err == nil && len(multiResp) > 0 {
			return multiResp[0].Embedding, nil
		}

		lastErr = fmt.Errorf("failed to parse response: %s", string(body))
		continue
	}

	return nil, fmt.Errorf("failed to generate embedding after 3 attempts: %w", lastErr)
}

// GenerateTableEmbedding generates embedding vector for table structure
func (h *HFEmbeddingService) GenerateTableEmbedding(table *models.Table) ([]float32, error) {
	// Construct textual representation of table structure
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

	return h.GenerateEmbedding(tableText)
}
