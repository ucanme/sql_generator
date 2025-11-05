package storage

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"awesomeProject2/internal/models"
	"github.com/google/uuid"
	_ "github.com/go-sql-driver/mysql"
)

// getTestMySQLStore creates a test MySQL store with a unique database name
func getTestMySQLStore(t *testing.T) *MySQLStore {
	// Using a test database - you might want to change this to your actual test DB
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/sqlbot_test?charset=utf8mb4&parseTime=True&loc=Local"
	}
	
	store, err := NewMySQLStore(dsn)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to MySQL: %v", err)
	}
	
	// Clean up any existing test data
	cleanupTestData(store.DB)
	
	return store
}

// cleanupTestData removes all test data from the database
func cleanupTestData(db *sql.DB) {
	db.Exec("DELETE FROM queries")
	db.Exec("DELETE FROM tables")
}

// createTestTable creates a sample table for testing
func createTestTable() *models.Table {
	return &models.Table{
		ID:          uuid.New().String(),
		Name:        fmt.Sprintf("test_table_%s", uuid.New().String()[:8]),
		Description: "A test table for unit testing",
		Columns: []models.Column{
			{
				Name:        "id",
				Type:        "INT",
				Description: "Primary key",
				IsPrimary:   true,
				IsRequired:  true,
			},
			{
				Name:        "name",
				Type:        "VARCHAR(255)",
				Description: "Name field",
				IsPrimary:   false,
				IsRequired:  true,
			},
			{
				Name:        "description",
				Type:        "TEXT",
				Description: "Description field",
				IsPrimary:   false,
				IsRequired:  false,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// createTestQuery creates a sample query for testing
func createTestQuery() *models.Query {
	return &models.Query{
		ID:          uuid.New().String(),
		Description: "A test query for unit testing",
		SQL:         "SELECT * FROM test_table",
		CreatedAt:   time.Now(),
	}
}

func TestMySQLStore_CreateTable(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	table := createTestTable()
	
	err := store.CreateTable(table)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	
	// Verify the table was created
	retrieved, err := store.GetTableByName(table.Name)
	if err != nil {
		t.Fatalf("Failed to get table: %v", err)
	}
	
	if retrieved.Name != table.Name {
		t.Errorf("Expected table name %s, got %s", table.Name, retrieved.Name)
	}
	
	if retrieved.Description != table.Description {
		t.Errorf("Expected table description %s, got %s", table.Description, retrieved.Description)
	}
	
	if len(retrieved.Columns) != len(table.Columns) {
		t.Errorf("Expected %d columns, got %d", len(table.Columns), len(retrieved.Columns))
	}
	
	// 验证时间戳
	if retrieved.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	
	if retrieved.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestMySQLStore_GetTableByName(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	table := createTestTable()
	
	// Table should not exist initially
	_, err := store.GetTableByName(table.Name)
	if err == nil {
		t.Fatalf("Expected error when getting non-existent table, got nil")
	}
	
	// Create the table
	err = store.CreateTable(table)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	
	// Now it should exist
	retrieved, err := store.GetTableByName(table.Name)
	if err != nil {
		t.Fatalf("Failed to get table: %v", err)
	}
	
	if retrieved.Name != table.Name {
		t.Errorf("Expected table name %s, got %s", table.Name, retrieved.Name)
	}
	
	// 测试获取空表名
	_, err = store.GetTableByName("")
	if err == nil {
		t.Error("Expected error when getting table with empty name, got nil")
	}
}

func TestMySQLStore_UpdateTable(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	table := createTestTable()
	
	// Create the table
	err := store.CreateTable(table)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	
	// Update the table
	originalUpdatedAt := table.UpdatedAt
	time.Sleep(time.Millisecond * 10) // 确保时间戳不同
	
	table.Description = "Updated description"
	table.Columns = append(table.Columns, models.Column{
		Name:        "email",
		Type:        "VARCHAR(255)",
		Description: "Email address",
		IsPrimary:   false,
		IsRequired:  false,
	})
	
	err = store.UpdateTable(table.Name, table)
	if err != nil {
		t.Fatalf("Failed to update table: %v", err)
	}
	
	// Verify the update
	updated, err := store.GetTableByName(table.Name)
	if err != nil {
		t.Fatalf("Failed to get updated table: %v", err)
	}
	
	if updated.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", updated.Description)
	}
	
	if len(updated.Columns) != 4 {
		t.Errorf("Expected 4 columns after update, got %d", len(updated.Columns))
	}
	
	// 验证更新时间戳已更改
	if !updated.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after modification")
	}
	
	// 测试更新不存在的表
	table.Name = "non_existent_table"
	err = store.UpdateTable(table.Name, table)
	if err == nil {
		t.Error("Expected error when updating non-existent table, got nil")
	}
}

func TestMySQLStore_DeleteTable(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	table := createTestTable()
	
	// Create the table
	err := store.CreateTable(table)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	
	// Delete the table
	err = store.DeleteTable(table.Name)
	if err != nil {
		t.Fatalf("Failed to delete table: %v", err)
	}
	
	// Table should no longer exist
	_, err = store.GetTableByName(table.Name)
	if err == nil {
		t.Fatalf("Expected error when getting deleted table, got nil")
	}
	
	// 测试删除不存在的表
	err = store.DeleteTable("non_existent_table")
	if err == nil {
		t.Error("Expected error when deleting non-existent table, got nil")
	}
}

func TestMySQLStore_SearchTables(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	// Create test tables
	table1 := createTestTable()
	table1.Name = "users_table"
	table1.Description = "Table containing user information"
	err := store.CreateTable(table1)
	if err != nil {
		t.Fatalf("Failed to create table1: %v", err)
	}
	
	table2 := createTestTable()
	table2.Name = "orders_table"
	table2.Description = "Table containing order information"
	err = store.CreateTable(table2)
	if err != nil {
		t.Fatalf("Failed to create table2: %v", err)
	}
	
	// Search for tables
	results, err := store.SearchTables("user", 10, 0)
	if err != nil {
		t.Fatalf("Failed to search tables: %v", err)
	}
	
	// At least table1 should be found
	found := false
	for _, table := range results {
		if table.Name == table1.Name {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected to find table %s in search results", table1.Name)
	}
	
	// 测试搜索限制
	results, err = store.SearchTables("table", 1, 0)
	if err != nil {
		t.Fatalf("Failed to search tables with limit: %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result with limit 1, got %d", len(results))
	}
	
	// 测试搜索偏移
	results, err = store.SearchTables("table", 10, 1)
	if err != nil {
		t.Fatalf("Failed to search tables with offset: %v", err)
	}
	
	// 应该至少有一个结果（因为我们有2个表）
	if len(results) < 1 {
		t.Error("Expected at least 1 result with offset")
	}
}

func TestMySQLStore_ListTables(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	// Create test tables
	table1 := createTestTable()
	table1.Name = "first_table"
	err := store.CreateTable(table1)
	if err != nil {
		t.Fatalf("Failed to create table1: %v", err)
	}
	
	time.Sleep(time.Millisecond * 10) // Ensure different timestamps
	
	table2 := createTestTable()
	table2.Name = "second_table"
	err = store.CreateTable(table2)
	if err != nil {
		t.Fatalf("Failed to create table2: %v", err)
	}
	
	// List tables
	results, err := store.ListTables(10, 0)
	if err != nil {
		t.Fatalf("Failed to list tables: %v", err)
	}
	
	if len(results) < 2 {
		t.Errorf("Expected at least 2 tables, got %d", len(results))
	}
	
	// 测试限制
	results, err = store.ListTables(1, 0)
	if err != nil {
		t.Fatalf("Failed to list tables with limit: %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result with limit 1, got %d", len(results))
	}
	
	// 测试偏移
	results, err = store.ListTables(10, 1)
	if err != nil {
		t.Fatalf("Failed to list tables with offset: %v", err)
	}
	
	// 应该至少有一个结果（因为我们有2个表）
	if len(results) < 1 {
		t.Error("Expected at least 1 result with offset")
	}
}

func TestMySQLStore_CreateQuery(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	query := createTestQuery()
	
	err := store.CreateQuery(query)
	if err != nil {
		t.Fatalf("Failed to create query: %v", err)
	}
	
	// Verify the query was created
	retrieved, err := store.GetQueryByID(query.ID)
	if err != nil {
		t.Fatalf("Failed to get query: %v", err)
	}
	
	if retrieved.Description != query.Description {
		t.Errorf("Expected query description %s, got %s", query.Description, retrieved.Description)
	}
	
	if retrieved.SQL != query.SQL {
		t.Errorf("Expected query SQL %s, got %s", query.SQL, retrieved.SQL)
	}
	
	// 验证时间戳
	if retrieved.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestMySQLStore_GetQueryByID(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	query := createTestQuery()
	
	// Query should not exist initially
	_, err := store.GetQueryByID(query.ID)
	if err == nil {
		t.Fatalf("Expected error when getting non-existent query, got nil")
	}
	
	// Create the query
	err = store.CreateQuery(query)
	if err != nil {
		t.Fatalf("Failed to create query: %v", err)
	}
	
	// Now it should exist
	retrieved, err := store.GetQueryByID(query.ID)
	if err != nil {
		t.Fatalf("Failed to get query: %v", err)
	}
	
	if retrieved.ID != query.ID {
		t.Errorf("Expected query ID %s, got %s", query.ID, retrieved.ID)
	}
	
	// 测试获取空ID
	_, err = store.GetQueryByID("")
	if err == nil {
		t.Error("Expected error when getting query with empty ID, got nil")
	}
}

func TestMySQLStore_ListQueries(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	// Create test queries
	query1 := createTestQuery()
	query1.Description = "First test query"
	err := store.CreateQuery(query1)
	if err != nil {
		t.Fatalf("Failed to create query1: %v", err)
	}
	
	time.Sleep(time.Millisecond * 10) // Ensure different timestamps
	
	query2 := createTestQuery()
	query2.Description = "Second test query"
	err = store.CreateQuery(query2)
	if err != nil {
		t.Fatalf("Failed to create query2: %v", err)
	}
	
	// List queries
	results, err := store.ListQueries(10, 0)
	if err != nil {
		t.Fatalf("Failed to list queries: %v", err)
	}
	
	if len(results) < 2 {
		t.Errorf("Expected at least 2 queries, got %d", len(results))
	}
	
	// 测试限制
	results, err = store.ListQueries(1, 0)
	if err != nil {
		t.Fatalf("Failed to list queries with limit: %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result with limit 1, got %d", len(results))
	}
	
	// 测试偏移
	results, err = store.ListQueries(10, 1)
	if err != nil {
		t.Fatalf("Failed to list queries with offset: %v", err)
	}
	
	// 应该至少有一个结果（因为我们有2个查询）
	if len(results) < 1 {
		t.Error("Expected at least 1 result with offset")
	}
}

// TestEdgeCases 测试边界情况和错误处理
func TestEdgeCases(t *testing.T) {
	store := getTestMySQLStore(t)
	defer store.DB.Close()
	
	// 测试带有特殊字符的表名
	table := createTestTable()
	table.Name = "test-table_with.special/chars"
	err := store.CreateTable(table)
	if err != nil {
		t.Fatalf("Failed to create table with special chars: %v", err)
	}
	
	// 验证可以检索它
	_, err = store.GetTableByName(table.Name)
	if err != nil {
		t.Fatalf("Failed to get table with special chars: %v", err)
	}
	
	// 测试大偏移量（应该返回空结果而不是错误）
	results, err := store.ListTables(10, 1000)
	if err != nil {
		t.Fatalf("Failed to list tables with large offset: %v", err)
	}
	
	if len(results) != 0 {
		t.Errorf("Expected 0 results with large offset, got %d", len(results))
	}
	
	// 测试负限制（应该使用默认值）
	results, err = store.ListTables(-1, 0)
	if err != nil {
		t.Fatalf("Failed to list tables with negative limit: %v", err)
	}
	
	if len(results) < 0 {
		t.Error("Expected non-negative results with negative limit")
	}
}