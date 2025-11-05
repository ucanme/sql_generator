package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"awesomeProject2/internal/models"
	"awesomeProject2/internal/storage"
	_ "github.com/go-sql-driver/mysql"
)

// getTestMySQLStore creates a test MySQL store with a unique database name
func getTestMySQLStore(t *testing.T) *storage.MySQLStore {
	// Using a test database - you might want to change this to your actual test DB
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/sqlbot_test?charset=utf8mb4&parseTime=True&loc=Local"
	}

	store, err := storage.NewMySQLStore(dsn)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to MySQL: %v", err)
	}

	// Clean up any existing test data
	cleanupTestData(t, store)

	return store
}

// cleanupTestData removes all test data from the database
func cleanupTestData(t *testing.T, store *storage.MySQLStore) {
	// Tables that need to be cleaned up (in proper order to avoid FK constraint issues)
	tables := []string{"audit_logs", "notifications", "time_tracking", "issue_comments", "issues",
		"documents", "task_comments", "tasks", "project_members", "projects", "user_roles",
		"roles", "user_departments", "departments", "users"}

	for _, tableName := range tables {
		_, err := store.DB.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
		if err != nil {
			t.Logf("Warning: Failed to clean up table %s: %v", tableName, err)
		}
	}

	// Also clean up queries
	_, err := store.DB.Exec("DELETE FROM queries")
	if err != nil {
		t.Logf("Warning: Failed to clean up queries table: %v", err)
	}
}

// generateSampleTables creates sample tables for testing
func generateSampleTables() []*models.Table {
	// This is a simplified version of GenerateSampleTables for integration tests
	tables := make([]*models.Table, 0, 3)

	now := time.Now()

	// Simple users table
	users := &models.Table{
		ID:          "1",
		Name:        "users",
		Description: "用户表",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "用户ID",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "username",
				Type:        "VARCHAR(50)",
				Description: "用户名",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "email",
				Type:        "VARCHAR(100)",
				Description: "邮箱",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, users)

	// Simple projects table
	projects := &models.Table{
		ID:          "2",
		Name:        "projects",
		Description: "项目表",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "项目ID",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "name",
				Type:        "VARCHAR(100)",
				Description: "项目名",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "项目描述",
				IsPrimary:   false,
				IsRequired:  false,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, projects)

	// Simple tasks table
	tasks := &models.Table{
		ID:          "3",
		Name:        "tasks",
		Description: "任务表",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "BIGINT",
				Description: "任务ID",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "project_id",
				Type:        "BIGINT",
				Description: "项目ID",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "title",
				Type:        "VARCHAR(200)",
				Description: "任务标题",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "status",
				Type:        "VARCHAR(20)",
				Description: "任务状态",
				IsPrimary:   false,
				IsRequired:  true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	tables = append(tables, tasks)

	return tables
}

// TestCompleteScenario tests the complete scenario with interconnected tables
func TestCompleteScenario(t *testing.T) {
	store := getTestMySQLStore(t)
	defer func() {
		// Clean up
		cleanupTestData(t, store)
	}()

	// Generate sample tables
	tables := generateSampleTables()

	// Insert all tables
	for _, table := range tables {
		err := store.CreateTable(table)
		if err != nil {
			t.Fatalf("Failed to create table %s: %v", table.Name, err)
		}
	}

	// Verify all tables were inserted
	for _, table := range tables {
		retrieved, err := store.GetTableByName(table.Name)
		if err != nil {
			t.Fatalf("Failed to retrieve table %s: %v", table.Name, err)
		}

		if retrieved.Name != table.Name {
			t.Errorf("Expected table name %s, got %s", table.Name, retrieved.Name)
		}

		if retrieved.Description != table.Description {
			t.Errorf("Expected table description %s, got %s", table.Description, retrieved.Description)
		}

		if len(retrieved.Columns) != len(table.Columns) {
			t.Errorf("Expected %d columns for table %s, got %d", len(table.Columns), table.Name, len(retrieved.Columns))
		}
	}

	// Test listing tables
	allTables, err := store.ListTables(20, 0)
	if err != nil {
		t.Fatalf("Failed to list tables: %v", err)
	}

	if len(allTables) != 3 {
		t.Errorf("Expected 3 tables in list, got %d", len(allTables))
	}

	// Test searching tables
	searchResults, err := store.SearchTables("用户", 10, 0)
	if err != nil {
		t.Fatalf("Failed to search tables: %v", err)
	}

	// Should find at least the users table
	foundUsers := false
	for _, table := range searchResults {
		if table.Name == "users" {
			foundUsers = true
			break
		}
	}

	if !foundUsers {
		t.Error("Expected to find users table in search results")
	}

	// Test updating a table
	usersTable, err := store.GetTableByName("users")
	if err != nil {
		t.Fatalf("Failed to get users table: %v", err)
	}

	// Add a new column
	newColumn := models.Column{
		Name:        "phone_number",
		Type:        "VARCHAR(20)",
		Description: "用户电话号码",
		IsPrimary:   false,
		IsRequired:  false,
	}
	usersTable.Columns = append(usersTable.Columns, newColumn)
	usersTable.UpdatedAt = time.Now()

	err = store.UpdateTable("users", usersTable)
	if err != nil {
		t.Fatalf("Failed to update users table: %v", err)
	}

	// Verify the update
	updatedUsersTable, err := store.GetTableByName("users")
	if err != nil {
		t.Fatalf("Failed to get updated users table: %v", err)
	}

	foundNewColumn := false
	for _, col := range updatedUsersTable.Columns {
		if col.Name == "phone_number" && col.Type == "VARCHAR(20)" {
			foundNewColumn = true
			break
		}
	}

	if !foundNewColumn {
		t.Error("Expected to find phone_number column in updated users table")
	}

	// Test deleting a table
	err = store.DeleteTable("tasks")
	if err != nil {
		t.Fatalf("Failed to delete tasks table: %v", err)
	}

	// Verify deletion
	_, err = store.GetTableByName("tasks")
	if err == nil {
		t.Error("Expected error when getting deleted tasks table, got nil")
	}

	// Verify we now have 2 tables
	allTables, err = store.ListTables(20, 0)
	if err != nil {
		t.Fatalf("Failed to list tables after deletion: %v", err)
	}

	if len(allTables) != 2 {
		t.Errorf("Expected 2 tables after deletion, got %d", len(allTables))
	}
}

// TestQueryOperations tests query operations with the sample data
func TestQueryOperations(t *testing.T) {
	store := getTestMySQLStore(t)
	defer func() {
		// Clean up
		cleanupTestData(t, store)
	}()

	// Generate and insert sample tables
	tables := generateSampleTables()
	for _, table := range tables {
		err := store.CreateTable(table)
		if err != nil {
			t.Fatalf("Failed to create table %s: %v", table.Name, err)
		}
	}

	// Create some sample queries
	queries := []*models.Query{
		{
			ID:          "1",
			Description: "查询所有用户信息",
			SQL:         "SELECT * FROM users",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Description: "查询指定项目的任务",
			SQL:         "SELECT * FROM tasks WHERE project_id = ?",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "3",
			Description: "查询项目信息",
			SQL:         "SELECT p.name, COUNT(t.id) as task_count FROM projects p LEFT JOIN tasks t ON p.id = t.project_id GROUP BY p.id",
			CreatedAt:   time.Now(),
		},
	}

	// Insert queries
	for _, query := range queries {
		err := store.CreateQuery(query)
		if err != nil {
			t.Fatalf("Failed to create query %s: %v", query.Description, err)
		}
	}

	// Retrieve and verify queries
	for _, expectedQuery := range queries {
		actualQuery, err := store.GetQueryByID(expectedQuery.ID)
		if err != nil {
			t.Fatalf("Failed to get query %s: %v", expectedQuery.ID, err)
		}

		if actualQuery.Description != expectedQuery.Description {
			t.Errorf("Expected query description %s, got %s", expectedQuery.Description, actualQuery.Description)
		}

		if actualQuery.SQL != expectedQuery.SQL {
			t.Errorf("Expected query SQL %s, got %s", expectedQuery.SQL, actualQuery.SQL)
		}
	}

	// List queries
	allQueries, err := store.ListQueries(10, 0)
	if err != nil {
		t.Fatalf("Failed to list queries: %v", err)
	}

	if len(allQueries) != 3 {
		t.Errorf("Expected 3 queries, got %d", len(allQueries))
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	store := getTestMySQLStore(t)
	defer func() {
		// Clean up
		cleanupTestData(t, store)
	}()

	// 测试空表名
	_, err := store.GetTableByName("")
	if err == nil {
		t.Error("Expected error for empty table name, got nil")
	}

	// 测试不存在的表
	_, err = store.GetTableByName("non_existent_table")
	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}

	// 测试空查询ID
	_, err = store.GetQueryByID("")
	if err == nil {
		t.Error("Expected error for empty query ID, got nil")
	}

	// 测试不存在的查询
	_, err = store.GetQueryByID("non_existent_query")
	if err == nil {
		t.Error("Expected error for non-existent query, got nil")
	}
}

// TestPagination 测试分页功能
func TestPagination(t *testing.T) {
	store := getTestMySQLStore(t)
	defer func() {
		// Clean up
		cleanupTestData(t, store)
	}()

	// Generate and insert sample tables
	tables := generateSampleTables()
	for _, table := range tables {
		err := store.CreateTable(table)
		if err != nil {
			t.Fatalf("Failed to create table %s: %v", table.Name, err)
		}
	}

	// 测试表分页
	page1, err := store.ListTables(2, 0)
	if err != nil {
		t.Fatalf("Failed to list tables page 1: %v", err)
	}

	if len(page1) != 2 {
		t.Errorf("Expected 2 tables in page 1, got %d", len(page1))
	}

	page2, err := store.ListTables(2, 2)
	if err != nil {
		t.Fatalf("Failed to list tables page 2: %v", err)
	}

	if len(page2) != 1 {
		t.Errorf("Expected 1 table in page 2, got %d", len(page2))
	}

	// 创建多个查询以测试查询分页
	for i := 0; i < 5; i++ {
		query := &models.Query{
			ID:          fmt.Sprintf("query_%d", i),
			Description: fmt.Sprintf("测试查询 %d", i),
			SQL:         "SELECT * FROM users",
			CreatedAt:   time.Now(),
		}
		err := store.CreateQuery(query)
		if err != nil {
			t.Fatalf("Failed to create query %d: %v", i, err)
		}
	}

	// 测试查询分页
	qPage1, err := store.ListQueries(3, 0)
	if err != nil {
		t.Fatalf("Failed to list queries page 1: %v", err)
	}

	if len(qPage1) != 3 {
		t.Errorf("Expected 3 queries in page 1, got %d", len(qPage1))
	}

	qPage2, err := store.ListQueries(3, 3)
	if err != nil {
		t.Fatalf("Failed to list queries page 2: %v", err)
	}

	if len(qPage2) != 2 {
		t.Errorf("Expected 2 queries in page 2, got %d", len(qPage2))
	}
}
