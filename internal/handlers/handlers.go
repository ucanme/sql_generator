package handlers

import (
	"net/http"
	"strconv"

	"awesomeProject2/internal/llm"
	"awesomeProject2/internal/models"
	"awesomeProject2/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler wraps all HTTP handlers
type Handler struct {
	store storage.Store
	llm   llm.Client
}

// NewHandler creates a new Handler
func NewHandler(store storage.Store, llmClient llm.Client) *Handler {
	return &Handler{
		store: store,
		llm:   llmClient,
	}
}

// RegisterRoutes registers all routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", h.HealthCheck)

	// Table routes
	tables := router.Group("/tables")
	{
		tables.POST("", h.CreateTable)
		tables.GET("", h.ListTables)
		tables.GET("/:name", h.GetTable)
		tables.PUT("/:name", h.UpdateTable)
		tables.DELETE("/:name", h.DeleteTable)
		tables.GET("/search/:keyword", h.SearchTables)
	}

	// Query routes
	queries := router.Group("/queries")
	{
		queries.POST("/generate", h.GenerateQuery)
		queries.GET("", h.ListQueries)
		queries.GET("/:id", h.GetQuery)
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   "SQL Query Bot is running",
	})
}

// CreateTable godoc
// @Summary Create a new table definition
// @Description Add a new table with its columns to the system
// @Tags tables
// @Accept json
// @Produce json
// @Param table body models.Table true "Table definition"
// @Success 201 {object} models.Table
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables [post]
func (h *Handler) CreateTable(c *gin.Context) {
	var table models.Table

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.CreateTable(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, table)
}

// GetTable godoc
// @Summary Get a table by name
// @Description Retrieve a table definition by its name
// @Tags tables
// @Produce json
// @Param name path string true "Table name"
// @Success 200 {object} models.Table
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables/{name} [get]
func (h *Handler) GetTable(c *gin.Context) {
	name := c.Param("name")

	table, err := h.store.GetTableByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// ListTables godoc
// @Summary List all tables
// @Description Get all table definitions with pagination
// @Tags tables
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {array} models.Table
// @Failure 500 {object} map[string]string
// @Router /tables [get]
func (h *Handler) ListTables(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tables, err := h.store.ListTables(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// SearchTables godoc
// @Summary Search tables by keyword
// @Description Search table definitions by keyword with pagination
// @Tags tables
// @Produce json
// @Param keyword path string true "Search keyword"
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {array} models.Table
// @Failure 500 {object} map[string]string
// @Router /tables/search/{keyword} [get]
func (h *Handler) SearchTables(c *gin.Context) {
	keyword := c.Param("keyword")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tables, err := h.store.SearchTables(keyword, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// UpdateTable godoc
// @Summary Update a table
// @Description Update a table definition by name
// @Tags tables
// @Accept json
// @Produce json
// @Param name path string true "Table name"
// @Param table body models.Table true "Updated table definition"
// @Success 200 {object} models.Table
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables/{name} [put]
func (h *Handler) UpdateTable(c *gin.Context) {
	name := c.Param("name")
	var table models.Table

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the name in the URL matches the name in the body
	table.Name = name

	if err := h.store.UpdateTable(name, &table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// DeleteTable godoc
// @Summary Delete a table
// @Description Remove a table definition by name
// @Tags tables
// @Param name path string true "Table name"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables/{name} [delete]
func (h *Handler) DeleteTable(c *gin.Context) {
	name := c.Param("name")

	if err := h.store.DeleteTable(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GenerateQuery godoc
// @Summary Generate SQL query
// @Description Generate SQL query based on natural language description
// @Tags queries
// @Accept json
// @Produce json
// @Param request body models.QueryRequest true "Query description"
// @Success 201 {object} models.Query
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /queries/generate [post]
func (h *Handler) GenerateQuery(c *gin.Context) {
	var req models.QueryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get relevant tables
	var tables []*models.Table
	var err error

	if len(req.TableNames) > 0 {
		// If table names are specified, get those tables
		tables, err = h.getSpecifiedTables(req.TableNames)
	} else {
		// Otherwise, search for relevant tables based on description
		tables, err = h.store.SearchTables(req.Description, 20, 0)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate SQL using LLM
	sql, err := h.llm.GenerateSQL(req.Description, tables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the query
	query := &models.Query{
		ID:          uuid.New().String(),
		Description: req.Description,
		SQL:         sql,
	}

	if err := h.store.CreateQuery(query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, query)
}

// getSpecifiedTables gets tables by their names
func (h *Handler) getSpecifiedTables(tableNames []string) ([]*models.Table, error) {
	var tables []*models.Table

	for _, name := range tableNames {
		table, err := h.store.GetTableByName(name)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// GetQuery godoc
// @Summary Get a generated query
// @Description Retrieve a previously generated query by ID
// @Tags queries
// @Produce json
// @Param id path string true "Query ID"
// @Success 200 {object} models.Query
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /queries/{id} [get]
func (h *Handler) GetQuery(c *gin.Context) {
	id := c.Param("id")

	query, err := h.store.GetQueryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, query)
}

// ListQueries godoc
// @Summary List all generated queries
// @Description Get all previously generated queries with pagination
// @Tags queries
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 100)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {array} models.Query
// @Failure 500 {object} map[string]string
// @Router /queries [get]
func (h *Handler) ListQueries(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	queries, err := h.store.ListQueries(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, queries)
}
