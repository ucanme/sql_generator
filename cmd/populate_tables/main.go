package main

import (
	"awesomeProject2/internal/config"
	"awesomeProject2/internal/models"
	"awesomeProject2/internal/storage"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create MySQL storage
	mysqlStore, err := storage.NewMySQLStore(cfg.MySQL.DSN)
	if err != nil {
		log.Fatalf("Failed to create MySQL store: %v", err)
	}
	defer mysqlStore.DB.Close()

	// Handle graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	<-sigChan
	// 	fmt.Println("\nReceived interrupt signal, shutting down...")
	// 	cancel()
	// }()

	// Check if a file path is provided as command line argument
	var tables []*models.Table
	if len(os.Args) > 1 {
		// Read tables from file
		filePath := os.Args[1]
		fmt.Printf("Reading tables from file: %s\n", filePath)

		// Determine file type by extension
		if strings.HasSuffix(strings.ToLower(filePath), ".sql") {
			tables, err = readTablesFromSQLFile(filePath)
		} else {
			tables, err = readTablesFromFile(filePath)
		}

		if err != nil {
			log.Fatalf("Failed to read tables from file: %v", err)
		}
	} else {
		// Generate sample tables
		fmt.Println("Generating sample tables is not supported when reading from file...")
		os.Exit(1)
	}

	// Generate SQL for the tables
	fmt.Println("\nGenerating SQL for tables...")
	for _, table := range tables {
		// Ensure table has an ID
		if table.ID == "" {
			table.ID = uuid.New().String()
		}

		// Ensure proper timestamps
		if table.CreatedAt.IsZero() {
			table.CreatedAt = time.Now()
		}
		table.UpdatedAt = time.Now()

		sql := generateTableSQL(table)
		fmt.Printf("SQL for table '%s':\n%s\n\n", table.Name, sql)
	}

	// Insert tables into database
	fmt.Println("Inserting tables into database...")
	insertedCount := 0
	for _, table := range tables {
		// Add a small delay to ensure unique timestamps
		time.Sleep(10 * time.Millisecond)

		// Ensure each table has a unique ID (reuse existing or generate new)
		if table.ID == "" {
			table.ID = uuid.New().String()
		}

		// Ensure proper timestamps
		if table.CreatedAt.IsZero() {
			table.CreatedAt = time.Now()
		}
		table.UpdatedAt = time.Now()

		err := mysqlStore.CreateTable(table)
		if err != nil {
			log.Printf("Failed to create table %s: %v", table.Name, err)
			continue
		}
		fmt.Printf("Successfully inserted table: %s\n", table.Name)
		insertedCount++
	}

	fmt.Printf("Successfully inserted %d tables into the database\n", insertedCount)

	// Verify tables were inserted
	fmt.Println("\nVerifying tables...")
	allTables, err := mysqlStore.ListTables(20, 0)
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	fmt.Printf("Found %d tables in the database:\n", len(allTables))
	for _, table := range allTables {
		fmt.Printf("- %s: %s\n", table.Name, table.Description)
	}

	fmt.Println("\nDatabase population completed successfully!")
}

// generateTableSQL generates SQL CREATE TABLE statement for a given table model
func generateTableSQL(table *models.Table) string {
	var sqlBuilder strings.Builder

	// Start CREATE TABLE statement
	sqlBuilder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", table.Name))

	// Add columns
	primaryKeys := []string{}
	for i, column := range table.Columns {
		// Add column definition
		sqlBuilder.WriteString(fmt.Sprintf("  `%s` %s", column.Name, column.Type))

		// Add NOT NULL constraint if required
		if column.IsRequired {
			sqlBuilder.WriteString(" NOT NULL")
		}

		// Add comment if exists
		if column.Description != "" {
			sqlBuilder.WriteString(fmt.Sprintf(" COMMENT '%s'", column.Description))
		}

		// Check if it's a primary key
		if column.IsPrimary {
			primaryKeys = append(primaryKeys, column.Name)
		}

		// Add comma if not the last column
		if i < len(table.Columns)-1 {
			sqlBuilder.WriteString(",\n")
		}
	}

	// Add primary key constraint if exists
	if len(primaryKeys) > 0 {
		sqlBuilder.WriteString(",\n")
		if len(primaryKeys) == 1 {
			sqlBuilder.WriteString(fmt.Sprintf("  PRIMARY KEY (`%s`)", primaryKeys[0]))
		} else {
			sqlBuilder.WriteString(fmt.Sprintf("  PRIMARY KEY (%s)", strings.Join(primaryKeys, "`, `")))
		}
	}

	// Close CREATE TABLE statement
	sqlBuilder.WriteString("\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;")

	// Add table comment
	if table.Description != "" {
		sqlBuilder.WriteString(fmt.Sprintf("\nALTER TABLE `%s` COMMENT='%s';", table.Name, table.Description))
	}

	return sqlBuilder.String()
}

// readTablesFromFile reads table structures from a JSON file
func readTablesFromFile(filePath string) ([]*models.Table, error) {
	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse JSON data
	var tables []*models.Table
	err = json.Unmarshal(data, &tables)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	// Assign IDs and timestamps if not present
	now := time.Now()
	for _, table := range tables {
		if table.ID == "" {
			table.ID = uuid.New().String()
		}
		if table.CreatedAt.IsZero() {
			table.CreatedAt = now
		}
		table.UpdatedAt = now
	}

	return tables, nil
}

// readTablesFromSQLFile reads table structures from a SQL file
func readTablesFromSQLFile(filePath string) ([]*models.Table, error) {
	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse SQL data to extract table structures
	content := string(data)

	// Find all CREATE TABLE statements
	createTableRegex := regexp.MustCompile(`(?is)CREATE TABLE(?: IF NOT EXISTS)?\s+[` + "`" + `"]?(\w+)[` + "`" + `"]?\s*\((.*?)(?:\)|,\s*PRIMARY KEY\s*\([^)]*\))`)
	matches := createTableRegex.FindAllStringSubmatch(content, -1)

	var tables []*models.Table
	now := time.Now()

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		tableName := match[1]
		columnsPart := match[2]

		table := &models.Table{
			ID:        uuid.New().String(),
			Name:      tableName,
			Columns:   []models.Column{},
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Parse columns
		columnLines := splitColumns(columnsPart)
		for _, line := range columnLines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Skip constraints
			if strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") ||
				strings.HasPrefix(strings.ToUpper(line), "FOREIGN KEY") ||
				strings.HasPrefix(strings.ToUpper(line), "CONSTRAINT") ||
				strings.HasPrefix(strings.ToUpper(line), "UNIQUE") ||
				strings.HasPrefix(strings.ToUpper(line), "KEY") {
				continue
			}

			// Parse column definition
			column := parseColumnDefinition(line)
			if column != nil {
				table.Columns = append(table.Columns, *column)
			}
		}

		// Try to extract table description from comments
		descriptionRegex := regexp.MustCompile(`(?i)COMMENT\s*=\s*['"](.*?)['"]`)
		descriptionMatches := descriptionRegex.FindStringSubmatch(content)
		if len(descriptionMatches) > 1 {
			table.Description = descriptionMatches[1]
		}

		tables = append(tables, table)
	}

	return tables, nil
}

// parseColumnDefinition parses a column definition line and returns a Column struct
func parseColumnDefinition(line string) *models.Column {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return nil
	}

	// Remove trailing comma
	trimmedLine = strings.TrimRight(trimmedLine, ",")

	// Use regex to parse column definition
	// Matches: column_name type [constraints]
	columnRegex := regexp.MustCompile(`(?i)^\s*[` + "`" + `"]?(\w+)[` + "`" + `"]?\s+(\w+(?:\([^)]*\))?)\s*(.*)$`)
	matches := columnRegex.FindStringSubmatch(trimmedLine)
	if len(matches) < 3 {
		return nil
	}

	name := matches[1]
	typePart := strings.ToUpper(matches[2])
	constraints := matches[3]

	// Extract column description from comments
	var description string
	commentRegex := regexp.MustCompile(`(?i)COMMENT\s*['"](.*?)['"]`)
	commentMatches := commentRegex.FindStringSubmatch(trimmedLine)
	if len(commentMatches) > 1 {
		description = commentMatches[1]
	}

	column := &models.Column{
		Name:        name,
		Type:        typePart,
		Description: description,
		IsRequired:  !strings.Contains(strings.ToUpper(constraints), "NULL") && !strings.Contains(strings.ToUpper(constraints), "DEFAULT"),
	}

	// Check if it's a primary key
	if strings.Contains(strings.ToUpper(constraints), "PRIMARY KEY") {
		column.IsPrimary = true
	}

	return column
}

// splitColumns splits the columns part of a CREATE TABLE statement into individual column definitions
func splitColumns(columnsPart string) []string {
	var columns []string
	var currentColumn strings.Builder
	var inParentheses int
	var inQuotes bool
	var quoteChar rune

	for _, r := range columnsPart {
		switch r {
		case '"', '`':
			if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else if quoteChar == r {
				inQuotes = false
			}
			currentColumn.WriteRune(r)
		case '(':
			if !inQuotes {
				inParentheses++
			}
			currentColumn.WriteRune(r)
		case ')':
			if !inQuotes {
				inParentheses--
			}
			currentColumn.WriteRune(r)
		case ',':
			if !inQuotes && inParentheses == 0 {
				columns = append(columns, currentColumn.String())
				currentColumn.Reset()
			} else {
				currentColumn.WriteRune(r)
			}
		default:
			currentColumn.WriteRune(r)
		}
	}

	// Add the last column if it exists
	if currentColumn.Len() > 0 {
		columns = append(columns, currentColumn.String())
	}

	return columns
}
