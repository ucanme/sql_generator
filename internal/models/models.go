package models

import "time"

// Table represents a database table structure
type Table struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name" binding:"required"`
	Description string    `json:"description" bson:"description"`
	Columns     []Column  `json:"columns" bson:"columns"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Column represents a column in a database table
type Column struct {
	Name        string `json:"name" bson:"name" binding:"required"`
	Type        string `json:"type" bson:"type" binding:"required"`
	Description string `json:"description" bson:"description"`
	IsPrimary   bool   `json:"is_primary" bson:"is_primary"`
	IsRequired  bool   `json:"is_required" bson:"is_required"`
}

// Query represents a generated SQL query
type Query struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Description string    `json:"description" bson:"description"`
	SQL         string    `json:"sql" bson:"sql"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// QueryRequest represents the request to generate a query
type QueryRequest struct {
	Description string   `json:"description" binding:"required"`
	TableNames  []string `json:"table_names,omitempty"`
}