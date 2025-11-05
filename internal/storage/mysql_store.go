package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"sql_generator/internal/models"
)

// MySQLStore implements Store interface with MySQL
type MySQLStore struct {
	DB *sql.DB
}

// NewMySQLStore creates a new MySQL storage
func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	// Create tables if they don't exist
	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Create indexes
	err = createIndexes(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &MySQLStore{
		DB: db,
	}, nil
}

// createTables creates the required tables if they don't exist
func createTables(db *sql.DB) error {
	tablesSQL := `
	CREATE TABLE IF NOT EXISTS tables (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		columns JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	queriesSQL := `
	CREATE TABLE IF NOT EXISTS queries (
		id VARCHAR(36) PRIMARY KEY,
		description TEXT,
		sql_text TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(tablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables table: %w", err)
	}

	_, err = db.Exec(queriesSQL)
	if err != nil {
		return fmt.Errorf("failed to create queries table: %w", err)
	}

	return nil
}

// createIndexes creates the required indexes
func createIndexes(db *sql.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_tables_name ON tables(name)",
		"CREATE INDEX IF NOT EXISTS idx_tables_created_at ON tables(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_queries_created_at ON queries(created_at)",
	}

	for _, indexSQL := range indexes {
		_, err := db.Exec(indexSQL)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// CreateTable saves a table definition
func (s *MySQLStore) CreateTable(table *models.Table) error {
	now := time.Now()
	table.CreatedAt = now
	table.UpdatedAt = now

	columnsJSON, err := json.Marshal(table.Columns)
	if err != nil {
		return fmt.Errorf("failed to marshal columns: %w", err)
	}

	_, err = s.DB.Exec(`
		INSERT INTO tables (id, name, description, columns, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, table.ID, table.Name, table.Description, columnsJSON, table.CreatedAt, table.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert table: %w", err)
	}

	return nil
}

// GetTableByName retrieves a table by name
func (s *MySQLStore) GetTableByName(name string) (*models.Table, error) {
	var table models.Table
	var columnsJSON []byte
	var createdAt, updatedAt time.Time

	err := s.DB.QueryRow(`
		SELECT id, name, description, columns, created_at, updated_at
		FROM tables
		WHERE name = ?
	`, name).Scan(&table.ID, &table.Name, &table.Description, &columnsJSON, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("table not found: %s", name)
		}
		return nil, fmt.Errorf("failed to find table: %w", err)
	}

	err = json.Unmarshal(columnsJSON, &table.Columns)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal columns: %w", err)
	}

	table.CreatedAt = createdAt
	table.UpdatedAt = updatedAt

	return &table, nil
}

// SearchTables searches tables by keywords
func (s *MySQLStore) SearchTables(keyword string, limit, offset int) ([]*models.Table, error) {
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Cap at 100 results
	}

	rows, err := s.DB.Query(`
		SELECT id, name, description, columns, created_at, updated_at
		FROM tables
		WHERE MATCH(description) AGAINST(? IN NATURAL LANGUAGE MODE)
		ORDER BY MATCH(description) AGAINST(? IN NATURAL LANGUAGE MODE) DESC
		LIMIT ? OFFSET ?
	`, keyword, keyword, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to search tables: %w", err)
	}
	defer rows.Close()

	var tables []*models.Table
	for rows.Next() {
		var table models.Table
		var columnsJSON []byte
		var createdAt, updatedAt time.Time

		err := rows.Scan(&table.ID, &table.Name, &table.Description, &columnsJSON, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}

		err = json.Unmarshal(columnsJSON, &table.Columns)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal columns: %w", err)
		}

		table.CreatedAt = createdAt
		table.UpdatedAt = updatedAt

		tables = append(tables, &table)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return tables, nil
}

// ListTables returns tables with pagination
func (s *MySQLStore) ListTables(limit, offset int) ([]*models.Table, error) {
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Cap at 100 results
	}

	rows, err := s.DB.Query(`
		SELECT id, name, description, columns, created_at, updated_at
		FROM tables
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []*models.Table
	for rows.Next() {
		var table models.Table
		var columnsJSON []byte
		var createdAt, updatedAt time.Time

		err := rows.Scan(&table.ID, &table.Name, &table.Description, &columnsJSON, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}

		err = json.Unmarshal(columnsJSON, &table.Columns)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal columns: %w", err)
		}

		table.CreatedAt = createdAt
		table.UpdatedAt = updatedAt

		tables = append(tables, &table)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return tables, nil
}

// UpdateTable updates a table by name
func (s *MySQLStore) UpdateTable(name string, table *models.Table) error {
	table.UpdatedAt = time.Now()

	columnsJSON, err := json.Marshal(table.Columns)
	if err != nil {
		return fmt.Errorf("failed to marshal columns: %w", err)
	}

	result, err := s.DB.Exec(`
		UPDATE tables
		SET description = ?, columns = ?, updated_at = ?
		WHERE name = ?
	`, table.Description, columnsJSON, table.UpdatedAt, name)

	if err != nil {
		return fmt.Errorf("failed to update table: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("table not found: %s", name)
	}

	return nil
}

// DeleteTable removes a table by name
func (s *MySQLStore) DeleteTable(name string) error {
	result, err := s.DB.Exec(`DELETE FROM tables WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("table not found: %s", name)
	}

	return nil
}

// CreateQuery saves a generated query
func (s *MySQLStore) CreateQuery(query *models.Query) error {
	query.CreatedAt = time.Now()

	_, err := s.DB.Exec(`
		INSERT INTO queries (id, description, sql_text, created_at)
		VALUES (?, ?, ?, ?)
	`, query.ID, query.Description, query.SQL, query.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert query: %w", err)
	}

	return nil
}

// GetQueryByID retrieves a query by ID
func (s *MySQLStore) GetQueryByID(id string) (*models.Query, error) {
	var query models.Query
	var createdAt time.Time

	err := s.DB.QueryRow(`
		SELECT id, description, sql_text, created_at
		FROM queries
		WHERE id = ?
	`, id).Scan(&query.ID, &query.Description, &query.SQL, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("query not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find query: %w", err)
	}

	query.CreatedAt = createdAt

	return &query, nil
}

// ListQueries returns all queries with pagination
func (s *MySQLStore) ListQueries(limit, offset int) ([]*models.Query, error) {
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Cap at 100 results
	}

	rows, err := s.DB.Query(`
		SELECT id, description, sql_text, created_at
		FROM queries
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list queries: %w", err)
	}
	defer rows.Close()

	var queries []*models.Query
	for rows.Next() {
		var query models.Query
		var createdAt time.Time

		err := rows.Scan(&query.ID, &query.Description, &query.SQL, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan query: %w", err)
		}

		query.CreatedAt = createdAt
		queries = append(queries, &query)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return queries, nil
}
